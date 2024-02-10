package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/netip"
	"os"
)

// Use this for now to go fast and declutter code
func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}

func SplitHTTPLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// read data one byte at a time and copy into token
	// if byte is \r, set a flag
	// if flag is set and byte it \n, return
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\r'); i >= 0 {

	}
}

func main() {
	// Listen for connections
	listenPort := must(netip.ParseAddrPort("127.0.0.1:4444"))
	listener := must(net.ListenTCP("tcp4", net.TCPAddrFromAddrPort(listenPort)))
	fmt.Println("Listening for connections...")
	clientConn := must(listener.Accept())
	fmt.Printf("Connection received from %v\n", clientConn.RemoteAddr())

	// Copy the client request to the upstream UNTIL \r\n\r\n
	// Then close the upstream and start listening for response (ignore request body for now)
	scanner := bufio.NewScanner(clientConn)
	scanner.Split()

	// Open upstream connection
	upstreamAddr := "127.0.0.1:8080"
	upstreamConn := must(net.Dial("tcp4", upstreamAddr))
	fmt.Println("Upstream connection established", upstreamConn.RemoteAddr())
	n, err := io.Copy(upstreamConn, clientConn)
	if err != nil {
		fmt.Printf("Error proxying: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Bytes sent upstream: %d", n)

	// Send request directly there

	// request := make([]byte, 4096)
	// //
	// idx := bytes.Index(request, []byte{'\r', '\n', '\r', '\n'})
	//
	// headers := make(map[string]string)
	// var requestLine string
	//
	// for i, headerPair := range bytes.Split(request[:idx], []byte("\r\n")) {
	// 	if i == 0 {
	// 		requestLine = string(headerPair)
	// 		continue
	// 	}
	// 	header := bytes.SplitN(headerPair, []byte(": "), 2)
	// 	headers[string(header[0])] = string(header[1])
	// }
	//
	// if _, ok := headers["Content-Length"]; ok {
	// 	fmt.Println("Content length but reading body not implemented")
	// }
	//
	// fmt.Println(requestLine)
	// fmt.Println(headers)
	// fmt.Println("Raw Body:", string(rawBody))
	//
	// conn.Write([]byte("HTTP/1.1 200\r\nContent-Type: text/plain\r\n\r\nGood\r\n"))
	clientConn.Close()
	listener.Close()
}
