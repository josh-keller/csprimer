package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

func sysCallShout(port int) {
	addr := unix.SockaddrInet4{
		Port: port,
		Addr: [4]byte{0, 0, 0, 0},
	}

	sfd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatalf("Socket: %v\n", err)
	}

	err = unix.Bind(sfd, &addr)
	if err != nil {
		log.Fatalf("Bind: %v", err)
	}

	buffer := make([]byte, 4096)
	for {
		n, from, err := unix.Recvfrom(sfd, buffer, 0)
		if err != nil {
			log.Fatalf("Recvfrom: %v", err)
		}
		message := string(buffer[0:n])

		if f, ok := from.(*unix.SockaddrInet4); ok {
			fmt.Printf("Message from %s:%d: %s", f.Addr, f.Port, message)
		} else {
			fmt.Printf("Message received: %s\n", message)
		}

		response := strings.ToUpper(message)

		err = unix.Sendto(sfd, []byte(response), 0, from)
		if err != nil {
			log.Fatalf("Sendto: %v", err)
		}
	}
}

func goShout(port int) {
	addrPort := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: port,
	}

	// syscalls: socket(), bind(), and listen()
	conn, err := net.ListenUDP("udp", addrPort)
	if err != nil {
		log.Fatalf("ListenUDP: %v", err)
	}
	fmt.Printf("Listening on port %d using Go interface\n", port)
	buffer := make([]byte, 4096)

	for {
		// syscall: recvfrom()
		n, from, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatalf("ReadFromUDP: %v", err)
		}
		message := string(buffer[0:n])
		fmt.Printf("Message from %s - %s", from.String(), message)

		// syscall: sendto()
		response := strings.ToUpper(message)
		conn.WriteToUDP([]byte(response), from)
		if err != nil {
			log.Fatalf("Write: %v", err)
		}
	}
}

func main() {
	portPtr := flag.Int("port", 2222, "Port number to use")
	unixPtr := flag.Bool("unix", false, "Call sysCallShout()")
	goPtr := flag.Bool("go", false, "Call goShout()")

	flag.Parse()

	if *unixPtr && *goPtr {
		fmt.Println("Error: Both -unix and -go flags cannot be used simultaneously.")
		os.Exit(1)
	}

	if *goPtr {
		goShout(*portPtr)
	} else {
		sysCallShout(*portPtr)
	}
}
