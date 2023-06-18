package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
)

type HTTPRequest struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    string
}

func ParseRequest(raw []byte) (HTTPRequest, error) {
	reqLineRE := regexp.MustCompile(`(?ms)^\s*(\w+)\s+([\w.\-\/]+)\s+(HTTP\/\d\.\d)\s*\r\n(.*)\r\n\r\n(.*)`)
	match := reqLineRE.FindSubmatch(raw)
	if len(match) == 0 {
		return HTTPRequest{}, errors.New("No match")
	}
	return HTTPRequest{
		Method:  string(match[1]),
		Path:    string(match[2]),
		Version: string(match[3]),
		Headers: parseHeaders(string(match[4])),
		Body:    string(match[5]),
	}, nil
}

func parseHeaders(raw string) map[string]string {
	fmt.Println("Raw:", raw)
	headerRE := regexp.MustCompile(`^\s*([\w\-]+):\s+(.*)$`)
	matches := headerRE.FindAllStringSubmatch(raw, 0)
	fmt.Println("Matches:", matches)
	headers := make(map[string]string)
	for _, match := range matches {
		headers[match[0]] = match[1]
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

		buffer := make([]byte, 4096)

		fmt.Print("reading...")
		n, err := conn.Read(buffer)
		fmt.Printf("read %d bytes\n", n)
		if err != nil {
			log.Print(err)
			break
		}
		log.Print("'", string(buffer[:n]), "'\n")
		req, err := ParseRequest(buffer)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%+v\n", req)
		}
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nOK"))
		if err = conn.Close(); err != nil {
			log.Print(err)
		}
	}
}
