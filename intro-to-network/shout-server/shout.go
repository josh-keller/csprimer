package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/sys/unix"
)

func main() {
	sfd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatalf("Socket: %v\n", err)
	}

	fmt.Printf("File descriptor: %d\n", sfd)

	addr := unix.SockaddrInet4{
		Port: 2222,
		Addr: [4]byte{0, 0, 0, 0},
	}

	err = unix.Bind(sfd, &addr)
	if err != nil {
		log.Fatalf("Bind: %v", err)
	}

	unix.Listen(sfd, 10)

	buffer := make([]byte, 4096)
	for {
		n, from, err := unix.Recvfrom(sfd, buffer, 0)
		if err != nil {
			log.Fatalf("Recvfrom: %v", err)
		}
		message := string(buffer[0:n])
		fmt.Printf("%s", message)
		response := strings.ToUpper(message)

		err = unix.Sendto(sfd, []byte(response), 0, from)
		if err != nil {
			log.Fatalf("Sendto: %v", err)
		}
	}
}
