package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// Goals:
// 1. Parse pcap file (use library if available)
// 2. Iterate through the packets, and display data

func main() {
	data := make(map[uint32][]byte)
	seqs := make([]uint32, 0)

	if handle, err := pcap.OpenOffline("./lossy.pcap"); err != nil {
		panic("Opening file")
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)
				if tcp.SrcPort == 80 {
					if _, found := data[tcp.Seq]; !found {
						data[tcp.Seq] = tcp.Payload
						seqs = append(seqs, tcp.Seq)
					}
					fmt.Printf("Inbound - Seq: %d | Ack: %d\n", tcp.Seq, tcp.Ack)
				}
			} else {
				fmt.Printf("Not TCP: %s\n", packet.String())
			}
		}
	}

	compiled := []byte{}
	slices.Sort(seqs)

	for _, s := range seqs {
		compiled = append(compiled, data[s]...)
	}

	// Find the Content-Length header
	dataHeaderIdx := bytes.Index(compiled, []byte("Content-Length: "))
	if dataHeaderIdx == -1 {
		panic("No content length found")
	}

	rawLength := []byte{}
	for idx := dataHeaderIdx + len("Content-Length: "); compiled[idx] != '\r' && compiled[idx] != '\n'; idx++ {
		rawLength = append(rawLength, compiled[idx])
	}

	fmt.Println("Raw length:", string(rawLength))

	// Find the two \r\n
	dataStart := bytes.Index(compiled, []byte("\r\n\r\n"))
	if dataStart == -1 {
		panic("dataStart")
	}

	dataStart += 4

	err := os.WriteFile("./output", compiled[dataStart:], 0600)
	if err != nil {
		panic("Writing output file")
	}
}
