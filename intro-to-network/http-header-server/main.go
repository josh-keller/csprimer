package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

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
	if req.Method != "GET" || req.Path != "/" {
		log.Printf("Responding not found to %s - (%s %s)", conn.RemoteAddr().String(), req.Method, req.Path)
		respondNotFound(conn)
	}

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

func respondNotFound(conn net.Conn) {
	respond(conn, "HTTP/1.1 404 Not Found", []byte("Not found"))
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
