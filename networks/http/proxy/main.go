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

func printErr(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func handleConn(clientConn net.Conn) {
	defer clientConn.Close()

	clientBuffer := make([]byte, 4096)
	n, err := clientConn.Read(clientBuffer)
	if err != nil {
		printErr("Error reading from client: %v", err)
		return
	}
	printErr("C --> P     U (%d bytes)", n)

	// For now, only check that it ends with \r\n\r\n.
	// TODO: handle requests that come in more than one packet
	if !bytes.ContainsAny(clientBuffer, "\r\n\r\n") {
		clientConn.Write(response500)
		printErr("Error reading from client: %v", err)
		panic("ERROR requests over 4k not implemented")
	}

	// Connect upstream and make request
	upstreamAddr := "127.0.0.1:8000"
	upstreamConn, err := net.Dial("tcp4", upstreamAddr)
	if err != nil {
		printErr("ERROR connecting upstream: %s", err)
		clientConn.Write(response502)
		return
	}
	defer upstreamConn.Close()

	upstreamRequest := clientBuffer

	n, err = upstreamConn.Write(upstreamRequest)
	if err != nil {
		clientConn.Write(response500)
		printErr("ERROR writing upstream: %s", err)
		return
	}
	printErr("C     P --> U (%d bytes)", n)

	// Loop to get upstream response
	totalBytesReturned := 0
	upstreamBuffer := make([]byte, 4096)
	for {
		// This blocks for a while at the end of each request.
		// Theory is that the connection doesn't get EOF because browser is using keep-alive
		// Eventually the browser closes and this call returns
		n, err = upstreamConn.Read(upstreamBuffer)
		totalBytesReturned += n
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			printErr("ERROR: reading from upstream: %s", err)
			clientConn.Write(response500)
			break
		}
		printErr("C     P <-- U (%d bytes)", n)
		_, err = clientConn.Write(upstreamBuffer[:n])
		if err != nil {
			printErr("ERROR: writing to client: %s", err)
			break
		}
		printErr("C <-- P     U (%d bytes)", n)
	}
	printErr("Client=%s, Upstream=%s, bytes=%d", clientConn.RemoteAddr(), upstreamConn.RemoteAddr(), totalBytesReturned)
}

func main() {
	port := 4444
	// Listen for connections
	listenAddr, _ := netip.ParseAddrPort(fmt.Sprintf("127.0.0.1:%d", port))
	listener, err := net.ListenTCP("tcp4", net.TCPAddrFromAddrPort(listenAddr))
	if err != nil {
		printErr("%v", err)
		os.Exit(1)
	}
	defer listener.Close()
	printErr("Listening for connections at %s...", listenAddr)
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			printErr("Error accepting connection: %s", err)
			continue
		}
		printErr("Connection received from %v", clientConn.RemoteAddr())

		go handleConn(clientConn)
	}
}
