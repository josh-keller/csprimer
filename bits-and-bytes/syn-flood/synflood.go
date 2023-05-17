package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
)

type PcapHeader struct {
	MagicNumber  uint32
	MajorVersion uint16
	MinorVersion uint16
	Reserved1    uint32
	Reserved2    uint32
	SnapLen      uint32
	FCS          uint16
	LinkType     uint16
}

type PacketHeader struct {
	TimestampSec   uint32
	TimestampMicro uint32
	CapturedLength uint32
	OriginalLength uint32
}

type IPPacket struct {
	VersionAndInfo     uint16
	Length             uint16
	Identification     uint16
	FlagsAndFragOffset uint16
	TTL                uint8
	Protocol           uint8
	Checksum           uint16
	SourceIP           [4]byte
	DestIP             [4]byte
}

type TCPSegment struct {
	SourcePort     uint16
	DestPort       uint16
	SequenceNumber uint32
	AckNumber      uint32
	DataOffset     uint8 // needs to be shifted
	Flags          uint8
	WindowSize     uint16
	Checksum       uint16
	UrgentPointer  uint16
}

type LatencyLog struct {
	SynReceived time.Time
	AckSent     time.Time
}

func main() {
	f, err := os.Open("./synflood.pcap")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)

	var pcapHeader PcapHeader
	binary.Read(r, binary.LittleEndian, &pcapHeader)
	fmt.Printf("Read PCAP Header: %+v\n\n", pcapHeader)

	var syn, ack int
	packetDataBuf := make([]byte, pcapHeader.SnapLen)
	var currPacketHeader PacketHeader
	var currPacket IPPacket
	var currSegment TCPSegment

	latencyLog := make(map[uint32]LatencyLog)

	for {
		err = binary.Read(r, binary.LittleEndian, &currPacketHeader)

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		// fmt.Printf("Pcap Header: %+v\n", currPacketHeader)
		io.ReadFull(r, packetDataBuf[:currPacketHeader.CapturedLength])
		packetBuf := bufio.NewReader(bytes.NewReader(packetDataBuf[:currPacketHeader.CapturedLength]))
		packetBuf.Discard(4) // Discard the link layer header

		binary.Read(packetBuf, binary.BigEndian, &currPacket)
		binary.Read(packetBuf, binary.BigEndian, &currSegment)
		// fmt.Printf("Packet: %+v\n", currPacket)
		// fmt.Printf("Segment: %+v\n", currSegment)

		if currSegment.Flags&0x10 == 0x10 {
			ack++
			latLog := latencyLog[currSegment.AckNumber]
			latLog.AckSent = time.Unix(int64(currPacketHeader.TimestampSec), 1000*int64(currPacketHeader.TimestampMicro))
			latencyLog[currSegment.AckNumber] = latLog
		} else if currSegment.Flags&0x02 == 0x02 {
			syn++
			latencyLog[currSegment.SequenceNumber] = LatencyLog{
				SynReceived: time.Unix(int64(currPacketHeader.TimestampSec), int64(currPacketHeader.TimestampMicro)*1000),
				AckSent:     time.Time{},
			}
		}
		// fmt.Print("\n-------------------\n")
	}

	fmt.Printf("Total Syn: %d, Total Ack: %d\n", syn, ack)

	var totalDuration time.Duration
	count := int64(0)

	for _, log := range latencyLog {
		if !log.AckSent.IsZero() {
			latency := log.AckSent.Sub(log.SynReceived)
			if latency < 0 {
				panic("Latency less than 0")
			}
			totalDuration += latency
			count++
		}
	}

	fmt.Println("Average latency: ", float64(totalDuration.Microseconds())/float64(count))
}
