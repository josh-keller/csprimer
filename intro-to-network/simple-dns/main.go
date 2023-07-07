package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"

	"golang.org/x/sys/unix"
)

const (
	// Bit 0
	DNS_QUERY = 0x0000
	DNS_REPLY = 0x8000
	// Bits 1-4
	OPCODE_QUERY   = 0x0000
	OPCODE_INVERSE = 0x0800
	OPCODE_STATUS  = 0x1000
	// Bit 5
	AUTHORITATIVE = 0x0400
	// Bit 6
	TRUNCATED = 0x0200
	// Bit 7
	RECURSION_DESIRED = 0x0100
	// Bit 8
	RECURSION_AVAIL = 0x0080
	// Bits 9-11 are Zero (reserved)
	// Bits 12-15
	RCODE_NOERROR         = 0x0000
	RCODE_FROMERR         = 0x0001
	RCODE_SERVFAIL        = 0x0002
	RCODE_NXDOMAIN        = 0x0003
	RCODE_NOT_IMPLEMENTED = 0x0004
	RCODE_REFUSED         = 0x0005
)

type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

func (h DNSHeader) RecursionAvailable() bool {
	return h.Flags&RECURSION_AVAIL != 0
}

func (h DNSHeader) RecursionDesired() bool {
	return h.Flags&RECURSION_DESIRED != 0
}

func (h DNSHeader) Query() bool {
	return h.Flags|DNS_QUERY == 0
}

func (h DNSHeader) Reply() bool {
	return h.Flags&DNS_REPLY != 0
}

func (h DNSHeader) OPCode() int {
	fmt.Printf("%x\n", h.Flags)
	fmt.Printf("%x\n", h.Flags&0x7800>>11)
	return int(h.Flags & 0x7800 >> 11)
}

func (h DNSHeader) ReturnCode() int {
	return int(h.Flags & 0x0007)
}

func (h *DNSHeader) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, h)
	return b.Bytes(), nil
}

func (h *DNSHeader) UnmarshalBinary(b []byte) error {
	r := bytes.NewReader(b)
	binary.Read(r, binary.BigEndian, h)
	return nil
}

var queryHeader = []byte{
	0xAB, 0xCD, // ID
	0x01, 0x00, // FLAGS
	0x00, 0x01, // Question Count
	0x00, 0x00, // Answer Count
	0x00, 0x00, // NS Count
	0x00, 0x00, // AR Count
}

var query = []byte{
	0x0a, 'j', 'o', 's', 'h', 'k', 'e', 'l', 'l', 'e', 'r',
	0x03, 'd', 'e', 'v',
	0x00,       // root
	0x00, 0x01, // A record
	0x00, 0x01, // Internet
}

func main() {
	addr := &unix.SockaddrInet4{
		Port: 53,
		Addr: [4]byte{8, 8, 8, 8},
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatalf("Socket: %v\n", err)
	}

	err = unix.Sendto(fd, append(queryHeader, query...), 0, addr)
	if err != nil {
		log.Fatalf("Send: %v\n", err)
	}

	buffer := make([]byte, 4096)
	n, _, err := unix.Recvfrom(fd, buffer, 0)
	if err != nil {
		log.Fatalf("Recv: %v\n", err)
	}

	fmt.Println(hex.EncodeToString(buffer[:n]), "\n----------------")
	fmt.Println(string(buffer[:n]))
}
