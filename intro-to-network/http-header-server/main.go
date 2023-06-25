package main

import (
	"encoding/json"
	"errors"
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
			log.Print(err)
		}

		handleRequest(conn)

	}
}

func handleRequest(conn *net.TCPConn) {
	buffer := make([]byte, 4096)

	_, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Problem reding request: %v\n", err)
		return
	} else {
		log.Printf("Request received from %s", conn.RemoteAddr().String())
	}

	req, err := ParseRequest(buffer)
	if err != nil {
		log.Printf("Problem parsing request: %v\n", err)
		return
	}

	conn.Write(constructResponse(req))
	if err = conn.Close(); err != nil {
		log.Printf("Problem writing response: %v\n", err)
	} else {
		log.Printf("Response returned to %s", conn.RemoteAddr().String())
	}

	return
}

func constructResponse(req *HTTPRequest) []byte {
	var b strings.Builder
	b.WriteString("HTTP/1.1 200 OK\r\n")
	b.WriteString("Content-Type: application/json\r\n\r\n")
	bodyBytes, err := json.Marshal(req.Headers)
	if err != nil {
		panic(err)
	}
	b.Write(bodyBytes)
	return []byte(b.String())
}
