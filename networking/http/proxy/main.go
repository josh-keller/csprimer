package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/netip"
	"os"
)

// TODO:
// - Send whole request to the upstream
// - Don't block waiting for EOF

// Use this for now to go fast and declutter code

/*
Receiving 400 from upstream.
- Tried changing Host
-
*/
func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}

var (
	response500 = []byte("HTTP/1.1 500 Internal Server Error\r\n\r\n")
	response502 = []byte("HTTP/1.1 502 Bad gateway\r\n\r\n")
)

func print_err(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func handleConn(clientConn net.Conn) {
	defer clientConn.Close()

	clientBuffer := make([]byte, 4096)
	n, err := clientConn.Read(clientBuffer)
	print_err("Received %d bytes from client", n)

	// For now, only check that it ends with \r\n\r\n.
	// TODO: handle requests that come in more than one packet
	if bytes.IndexAny(clientBuffer, "\r\n\r\n") == -1 {
		clientConn.Write(response500)
		panic("ERROR requests over 4k not implemented")
	}

	// Connect upstream and make request
	upstreamAddr := "127.0.0.1:8080"
	upstreamConn, err := net.Dial("tcp4", upstreamAddr)
	if err != nil {
		print_err("ERROR connecting upstream: %s", err)
		clientConn.Write(response502)
		return
	}
	defer upstreamConn.Close()

	// Get just the request line and write it upstream
	// TODO: figure out why upstream wasn't working with full request
	// requestLine := bytes.Split(clientBuffer, []byte("\r\n"))[0]
	// requestLine = append(requestLine, []byte("\r\n\r\n")...)
	// re := regexp.MustCompile(`Host: .*`)
	// print_err("Raw request:\n%s", clientBuffer)
	// upstreamRequest := re.ReplaceAllLiteral(clientBuffer, []byte("Host: "+upstreamConn.RemoteAddr().String()))
	// print_err("Rewritten request:\n%s", upstreamRequest)
	upstreamRequest := clientBuffer

	n, err = upstreamConn.Write(upstreamRequest)
	if err != nil {
		clientConn.Write(response500)
		print_err("ERROR writing upstream: %s", err)
		return
	}

	// Loop to get upstream response
	totalBytesReturned := 0
	upstreamBuffer := make([]byte, 4096)
	for {
		// This blocks for a while at the end of each request.
		// Theory is that the connection doesn't get EOF because browser is using keep-alive
		// Eventually the browser closes and this call returns
		n, err = upstreamConn.Read(upstreamBuffer)
		print_err("Upstream buffer: %s", upstreamBuffer)
		totalBytesReturned += n
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			print_err("ERROR: reading from upstream: %s", err)
			clientConn.Write(response500)
			break
		}
		_, err = clientConn.Write(upstreamBuffer)
		if err != nil {
			print_err("ERROR: writing to client: %s", err)
			break
		}
	}
	print_err("%s <-- %s: %d bytes transferred", clientConn.RemoteAddr(), upstreamConn.RemoteAddr(), totalBytesReturned)
}

func main() {
	port := 4444
	// Listen for connections
	listenAddr := must(netip.ParseAddrPort(fmt.Sprintf("127.0.0.1:%d", port)))
	listener := must(net.ListenTCP("tcp4", net.TCPAddrFromAddrPort(listenAddr)))
	print_err("Listening for connections at %s...", listenAddr)
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			print_err("Error accepting connection: %s", err)
			continue
		}
		print_err("Connection received from %v", clientConn.RemoteAddr())

		go handleConn(clientConn)
	}

	listener.Close()
}
