package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type Header struct {
	Signature string
	Size      uint32
	Offset    uint32
}

type InfoHeader struct {
	Size         uint32
	Width        uint32
	Height       uint32
	Planes       uint16
	BitsPerPixel uint16
	Compression  uint32
	ImageSize    uint32
	ColorsUsed   uint32
}

func main() {
	file, _ := os.Open("./input/teapot.bmp")

	header := ReadHeader(file)
	infoHeader := ReadInfoHeader(file)

	fmt.Printf("Header: %+v\n", header)
	fmt.Printf("Info Header: %+v\n", infoHeader)

	rawHeaders := make([]byte, header.Offset)
	rawPic := make([]byte, infoHeader.ImageSize)

	file.Read(rawHeaders)
	file.Read(rawPic)

	width := infoHeader.Width
	height := infoHeader.Height

	rotatedBuffer := make([]byte, 3*width*height)

	for r := uint32(0); r < height; r++ {
		for c := uint32(0); c < width; c++ {
			cnew := r
			rnew := height - 1 - c

			oldIdx := r*width*3 + 3*c
			newIdx := rnew*width*3 + 3*cnew

			for i := uint32(0); i < 3; i++ {
				rotatedBuffer[newIdx+i] = rawPic[oldIdx+i]
			}
		}
	}

	outFile, _ := os.Create("output.bmp")
	outFile.Write(rawHeaders)
	outFile.Write(rotatedBuffer)
	outFile.Sync()
	outFile.Close()
}

func ReadHeader(file *os.File) Header {
	header := make([]byte, 14)
	file.ReadAt(header, 0)
	sig := header[0:2]
	size := binary.LittleEndian.Uint32(header[2:6])
	offset := binary.LittleEndian.Uint32(header[10:14])
	return Header{string(sig), size, offset}
}

func ReadInfoHeader(file *os.File) InfoHeader {
	s := make([]byte, 4)
	file.ReadAt(s, 0x000E)
	size := binary.LittleEndian.Uint32(s)

	infoHeader := InfoHeader{Size: size}

	buffer := make([]byte, size-4)
	file.ReadAt(buffer, 0x0012)
	infoHeader.Width = binary.LittleEndian.Uint32(buffer[0:4])
	infoHeader.Height = binary.LittleEndian.Uint32(buffer[4:8])
	infoHeader.Planes = binary.LittleEndian.Uint16(buffer[8:10])
	infoHeader.BitsPerPixel = binary.LittleEndian.Uint16(buffer[10:12])
	infoHeader.Compression = binary.LittleEndian.Uint32(buffer[12:16])
	infoHeader.ImageSize = binary.LittleEndian.Uint32(buffer[16:20])
	infoHeader.ColorsUsed = binary.LittleEndian.Uint32(buffer[28:32])

	return infoHeader
}
