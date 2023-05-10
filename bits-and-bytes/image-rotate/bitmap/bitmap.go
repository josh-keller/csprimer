package bitmap

import (
	"encoding/binary"
	"io"
)

type Header struct {
	Signature [2]byte
	FileSize  uint32
	Reserved  uint32
	Offset    uint32
}

type InfoHeader struct {
	InfoSize     uint32
	Width        uint32
	Height       uint32
	Planes       uint16
	BitsPerPixel uint16
	Compression  uint32
	ImageSize    uint32
	XpixelsPerM  uint32
	YpixelsPerM  uint32
	ColorsUsed   uint32
	RawInfo      []byte
}

type Pixel struct {
	Green uint8
	Blue  uint8
	Red   uint8
}

type Bitmap struct {
	Header
	InfoHeader
	Image []byte
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

	bmp.Image = readPixels(r, bmp.ImageSize)
	return bmp
}

func readPixels(r io.Reader, size uint32) []byte {
	pixels := make([]byte, size)
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
	infoHeader := InfoHeader{InfoSize: size}

	buffer := make([]byte, size-4)
	r.Read(buffer)
	infoHeader.Width = binary.LittleEndian.Uint32(buffer[0:4])
	infoHeader.Height = binary.LittleEndian.Uint32(buffer[4:8])
	infoHeader.Planes = binary.LittleEndian.Uint16(buffer[8:10])
	infoHeader.BitsPerPixel = binary.LittleEndian.Uint16(buffer[10:12])
	infoHeader.Compression = binary.LittleEndian.Uint32(buffer[12:16])
	infoHeader.ImageSize = binary.LittleEndian.Uint32(buffer[16:20])
	infoHeader.ColorsUsed = binary.LittleEndian.Uint32(buffer[28:32])
	infoHeader.RawInfo = append(binary.LittleEndian.AppendUint32(infoHeader.RawInfo, size), buffer...)

	return infoHeader
}

func Rotate(bmp Bitmap) Bitmap {
	var rotated Bitmap
	rotated.Header = bmp.Header
	rotated.InfoHeader = bmp.InfoHeader
	rotated.Width, rotated.Height = rotated.Height, rotated.Width
	for i := 0; i < 4; i++ {
		rotated.RawInfo[4+i], rotated.RawInfo[8+i] = rotated.RawInfo[8+i], rotated.RawInfo[4+i]
	}

	width := bmp.Width
	height := bmp.Height
	newWidth := rotated.Width
	newHeight := rotated.Height
	oldPad := (4 - ((width * 3) % 4)) % 4
	newPad := (4 - ((newWidth * 3) % 4)) % 4

	rotated.InfoHeader.ImageSize = (3*newWidth + newPad) * newHeight
	rotated.Image = make([]byte, rotated.ImageSize)

	for r := uint32(0); r < height; r++ {
		for c := uint32(0); c < width; c++ {
			cnew := r
			rnew := newHeight - 1 - c

			oldIdx := r*(width*3+oldPad) + 3*c
			newIdx := rnew*(newWidth*3+newPad) + 3*cnew
			rotated.Image[newIdx] = bmp.Image[oldIdx]
			rotated.Image[newIdx+1] = bmp.Image[oldIdx+1]
			rotated.Image[newIdx+2] = bmp.Image[oldIdx+2]
		}
	}

	return rotated
}
