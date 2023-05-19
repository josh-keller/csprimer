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

type timeStampPair struct {
	SynReceived time.Time
	AckSent     time.Time
}

func (s TCPSegment) Syn() bool {
	return s.Flags&0x02 == 0x02
}

func (s TCPSegment) Ack() bool {
	return s.Flags&0x10 == 0x10
}

type LatencyLog struct {
	log map[uint32]timeStampPair
}

func (l *LatencyLog) Add(ph PacketHeader, seg TCPSegment) {
	if seg.Ack() {
		pair := l.log[seg.AckNumber-1]
		pair.AckSent = time.Unix(int64(ph.TimestampSec), 1000*int64(ph.TimestampMicro))
		l.log[seg.AckNumber-1] = pair
		return
	}
	if seg.Syn() {
		l.log[seg.SequenceNumber] = timeStampPair{
			SynReceived: time.Unix(int64(ph.TimestampSec), int64(ph.TimestampMicro)*1000),
			AckSent:     time.Time{},
		}
		return
	}
}

func (l *LatencyLog) Average() float64 {
	var totalDuration time.Duration
	count := int64(0)

	for _, log := range l.log {
		if !log.AckSent.IsZero() {
			latency := log.AckSent.Sub(log.SynReceived)
			if latency < 0 {
				panic("Latency less than 0")
			}
			totalDuration += latency
			count++
		}
	}

	return float64(totalDuration.Seconds()) / float64(count)
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

	packetDataBuf := make([]byte, pcapHeader.SnapLen)
	var synCount, ackCount int
	var currPacketHeader PacketHeader
	var currPacket IPPacket
	var currSegment TCPSegment
	latencyLog := LatencyLog{log: make(map[uint32]timeStampPair)}

	for {
		// Read the packet header
		err = binary.Read(r, binary.LittleEndian, &currPacketHeader)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		// Read in all data for packet
		io.ReadFull(r, packetDataBuf[:currPacketHeader.CapturedLength])
		packetBuf := bufio.NewReader(bytes.NewReader(packetDataBuf[:currPacketHeader.CapturedLength]))

		// Parse into structs
		packetBuf.Discard(4) // Discard the link layer header
		binary.Read(packetBuf, binary.BigEndian, &currPacket)
		binary.Read(packetBuf, binary.BigEndian, &currSegment)

		// Check for flags
		if currSegment.Ack() {
			ackCount++
			latencyLog.Add(currPacketHeader, currSegment)
		} else if currSegment.Syn() {
			synCount++
			latencyLog.Add(currPacketHeader, currSegment)
		}
	}

	fmt.Printf("Total Syn: %d, Total Ack: %d\n", synCount, ackCount)
	fmt.Printf("Average latency: %.03fs\n", latencyLog.Average())
}
