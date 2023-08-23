package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/netip"
)

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}

func main() {
	port := must(netip.ParseAddrPort("127.0.0.1:4444"))
	listener := must(net.ListenTCP("tcp4", net.TCPAddrFromAddrPort(port)))
	conn := must(listener.Accept())
	request := make([]byte, 4096)

	n := must(conn.Read(request))
	scanner := bufio.NewScanner(conn)
	idx := bytes.Index(request, []byte{'\r', '\n', '\r', '\n'})

	headers := make(map[string]string)
	var requestLine string

	for i, headerPair := range bytes.Split(request[:idx], []byte("\r\n")) {
		if i == 0 {
			requestLine = string(headerPair)
			continue
		}
		header := bytes.SplitN(headerPair, []byte(": "), 2)
		headers[string(header[0])] = string(header[1])
	}

	if _, ok := headers["Content-Length"]; ok {
		fmt.Println("Content length but reading body not implemented")
	}

	fmt.Println(requestLine)
	fmt.Println(headers)
	fmt.Println("Raw Body:", string(rawBody))

	conn.Write([]byte("HTTP/1.1 200\r\nContent-Type: text/plain\r\n\r\nGood\r\n"))
	conn.Close()
	listener.Close()
}
