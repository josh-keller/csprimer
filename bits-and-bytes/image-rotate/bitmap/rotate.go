package bitmap

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type Header struct {
	Signature [2]byte
	FileSize  uint32
	Reserved  uint32
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
	XpixelsPerM  uint32
	YpixelsPerM  uint32
	ColorsUsed   uint32
	Raw          []byte
}

type Pixel struct {
	Green uint8
	Blue  uint8
	Red   uint8
}

type Bitmap struct {
	Header     Header
	InfoHeader InfoHeader
	Image      []Pixel
}

// TODO: Go back to old Image implementaiton: just bytes, not Pixel struct
// This will help the padding
// Calculate padding based on width before reading, make sure to account for padding while reading
// Calculate padding when rotating too
// Pull this out into function of r, c

func NewFromReader(r io.Reader) Bitmap {
	var bmp Bitmap

	bmp.Header = readHeader(r)
	sizeBuffer := make([]byte, 4)
	r.Read(sizeBuffer)
	ihSize := binary.LittleEndian.Uint32(sizeBuffer)
	bmp.InfoHeader = readInfoHeader(r, ihSize)

	bmp.Image = readPixels(r, bmp.InfoHeader.ImageSize)
	return bmp
}

func readPixels(r io.Reader, size uint32) []Pixel {
	pixels := make([]Pixel, size)
	for i := uint32(0); i < size; i++ {
		binary.Read(r, binary.BigEndian, &pixels[i])
	}

	return pixels
}

func readHeader(r io.Reader) Header {
	h := Header{}
	binary.Read(r, binary.LittleEndian, &h)
	return h
}

func readInfoHeader(r io.Reader, size uint32) InfoHeader {
	infoHeader := InfoHeader{Size: size}

	buffer := make([]byte, size-4)
	r.Read(buffer)
	infoHeader.Width = binary.LittleEndian.Uint32(buffer[0:4])
	infoHeader.Height = binary.LittleEndian.Uint32(buffer[4:8])
	infoHeader.Planes = binary.LittleEndian.Uint16(buffer[8:10])
	infoHeader.BitsPerPixel = binary.LittleEndian.Uint16(buffer[10:12])
	infoHeader.Compression = binary.LittleEndian.Uint32(buffer[12:16])
	infoHeader.ImageSize = binary.LittleEndian.Uint32(buffer[16:20])
	infoHeader.ColorsUsed = binary.LittleEndian.Uint32(buffer[28:32])
	infoHeader.Raw = append(binary.LittleEndian.AppendUint32(infoHeader.Raw, size), buffer...)

	return infoHeader
}

func Rotate(bmp Bitmap) Bitmap {
	var rotated Bitmap
	rotated.Header = bmp.Header
	rotated.InfoHeader = bmp.InfoHeader
	rotated.InfoHeader.Width, rotated.InfoHeader.Height = rotated.InfoHeader.Height, rotated.InfoHeader.Width
	for i := 0; i < 4; i++ {
		rotated.InfoHeader.Raw[4+i], rotated.InfoHeader.Raw[8+i] = rotated.InfoHeader.Raw[8+i], rotated.InfoHeader.Raw[4+i]
	}

	width := bmp.InfoHeader.Width
	height := bmp.InfoHeader.Height
	newWidth := rotated.InfoHeader.Width
	newHeight := rotated.InfoHeader.Height
	oldPad := 4 - ((width * 3) % 4)
	newPad := 4 - ((newWidth * 3) % 4)

	rotated.InfoHeader.ImageSize = (newWidth + newPad) * newHeight
	rotated.Image = make([]Pixel, rotated.InfoHeader.ImageSize)

	for r := uint32(0); r < height; r++ {
		for c := uint32(0); c < width; c++ {
			cnew := r
			rnew := newHeight - 1 - c

			oldIdx := r*(width+oldPad) + c
			newIdx := rnew*(newWidth+newPad) + cnew
			// Problem: padding should be in bytes, but it's in Pixels
			if c == 0 || c == width-1 {
				fmt.Fprintf(os.Stderr, "old: %x,%x,%x new: %x,%x,%x\n", r, c, oldIdx, rnew, cnew, newIdx)
			}
			rotated.Image[newIdx] = bmp.Image[oldIdx]
		}
	}

	return rotated
}
