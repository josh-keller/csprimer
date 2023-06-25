package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

type HTTPRequest struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    string
}

func ParseRequest(raw []byte) (*HTTPRequest, error) {
	reqLineRE := regexp.MustCompile(`(?ms)^\s*(\w+)\s+([\w.\-\/]+)\s+(HTTP\/\d\.\d)\s*\r\n(.*)\r\n\r\n(.*)`)
	match := reqLineRE.FindSubmatch(raw)
	if len(match) == 0 {
		return nil, errors.New("No match")
	}
	return &HTTPRequest{
		Method:  string(match[1]),
		Path:    string(match[2]),
		Version: string(match[3]),
		Headers: parseHeaders(string(match[4])),
		Body:    string(match[5]),
	}, nil
}

func parseHeaders(raw string) map[string]string {
	headerRE := regexp.MustCompile(`(?mi)^\s*([\w\-]+):\s+(.*)$`)
	matches := headerRE.FindAllStringSubmatch(raw, -1)
	headers := make(map[string]string)
	for _, match := range matches {
		headers[match[1]] = strings.TrimSpace(match[2])
	}

	return headers
}

func main() {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 8080,
	})

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Problem accepting connection: %v\n", err)
			continue
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn *net.TCPConn) {
	defer conn.Close()
	buffer := make([]byte, 4096)

	_, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Problem reading request: %v\n", err)
		respondBadRequest(conn, err)
		return
	}
	req, err := ParseRequest(buffer)
	if err != nil {
		log.Printf("Problem parsing request: %v\n", err)
		respondBadRequest(conn, err)
		return
	}

	log.Printf("Request received from %s - %s %s", conn.RemoteAddr().String(), req.Method, req.Path)

	body, err := json.Marshal(req.Headers)
	if err != nil {
		log.Printf("Problem constructing response body: %v\n", err)
		respondInternalError(conn, err)
	}

	bytesWritten, err := respondOK(conn, body)

	if err != nil {
		log.Printf("Problem writing response: %v\n", err)
	} else {
		log.Printf("%d bytes returned to %s", bytesWritten, conn.RemoteAddr().String())
	}

	return
}

func respond(conn net.Conn, statusLine string, body []byte) (int, error) {
	n, err := conn.Write([]byte(statusLine +
		"\r\nContent-Type: application/json\r\n" +
		fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))),
	)
	if err != nil {
		return n, err
	}

	b, err := conn.Write(body)
	if err != nil {
		return n, err
	}

	return b + n, nil
}

func respondOK(conn net.Conn, body []byte) (int, error) {
	return respond(conn, "HTTP/1.1 200 OK", body)
}

func respondBadRequest(conn net.Conn, err error) {
	body, _ := json.Marshal(struct{ Error string }{Error: err.Error()})
	respond(conn, "HTTP/1.1 400 Bad Request", body)
	return
}

func respondInternalError(conn net.Conn, err error) {
	body, _ := json.Marshal(struct{ Error string }{Error: err.Error()})
	respond(conn, "HTTP/1.1 500 Internal Server Error", body)
	return
}
