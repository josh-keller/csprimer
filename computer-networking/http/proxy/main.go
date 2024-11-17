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

// Use this for now to go fast and declutter code
func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}

// func SplitHTTPLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
// 	// read data one byte at a time and copy into token
// 	// if byte is \r, set a flag
// 	// if flag is set and byte it \n, return
// 	if atEOF && len(data) == 0 {
// 		return 0, nil, nil
// 	}
//
// 	if i := bytes.IndexByte(data, '\r'); i >= 0 {
//
// 	}
// }

func main() {
	// Listen for connections
	listenPort := must(netip.ParseAddrPort("127.0.0.1:4444"))
	listener := must(net.ListenTCP("tcp4", net.TCPAddrFromAddrPort(listenPort)))
	fmt.Fprintln(os.Stderr, "Listening for connections...")
	for {
		clientConn := must(listener.Accept())
		fmt.Fprintf(os.Stderr, "Connection received from %v\n", clientConn.RemoteAddr())
		// Read into clientBuffer
		clientBuffer := make([]byte, 4096)
		clientConn.Read(clientBuffer)

		fmt.Fprintf(os.Stderr, "Received data: '%s'\n", clientBuffer)
		// For now, only check that it ends with \r\n\r\n.
		// TODO: handle requests that come in more than one packet
		if bytes.IndexAny(clientBuffer, "\r\n\r\n") == -1 {
			panic("Not implemented: multi-packet requests")
		}

		// Connect upstream and make request
		upstreamAddr := "127.0.0.1:8080"
		upstreamConn := must(net.Dial("tcp4", upstreamAddr))
		fmt.Fprintf(os.Stderr, "Upstream connection established: %s\n", upstreamConn.RemoteAddr())
		requestLine := bytes.Split(clientBuffer, []byte("\r\n"))[0]
		requestLine = append(requestLine, []byte("\r\n\r\n")...)
		n, err := upstreamConn.Write(requestLine)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing upstream: %s", err)
		}
		// Get just the request line and write it upstream
		fmt.Fprintf(os.Stderr, "Wrote %d bytes upstream\n", n)
		// TODO: handle error cases
		upstreamBuffer := make([]byte, 4096)
		for {
			n, err = upstreamConn.Read(upstreamBuffer)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: reading from upstream: %s", err)
				clientConn.Write([]byte("HTTP/1.1 500 Server Error\r\n\r\n"))
				break
			}
			fmt.Fprintf(os.Stderr, "Received %d bytes from upstream:\n%s\n\n", n, upstreamBuffer)
			_, err = clientConn.Write(upstreamBuffer)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: writing to client: %s", err)
			}
		}

		upstreamConn.Close()
		clientConn.Close()
		fmt.Fprintf(os.Stderr, "\n---------------------\n")
	}

	// Read into a buffer

	// Copy the client request to the upstream UNTIL \r\n\r\n
	// Then close the upstream and start listening for response (ignore request body for now)
	//
	// // Open upstream connection
	// n, err := io.Copy(upstreamConn, clientConn)
	// if err != nil {
	// 	fmt.Printf("Error proxying: %v\n", err)
	// 	os.Exit(1)
	// }
	//
	// fmt.Printf("Bytes sent upstream: %d", n)
	//
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
	listener.Close()
}
