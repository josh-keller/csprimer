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

type FP2DOT30 uint32

type InfoHeader struct {
	InfoSize uint32
	Info
}

type Info struct {
	Width           int32
	Height          int32
	Planes          uint16
	BitsPerPixel    uint16
	Compression     uint32
	ImageSize       uint32
	XpixelsPerM     int32
	YpixelsPerM     int32
	ColorsUsed      uint32
	ImportantColors uint32
	// RedMask         uint32
	// GreenMask       uint32
	// BlueMask        uint32
	// AlphaMask       uint32
	// CSType          uint32
	// Endpoints       [3]uint32
	// GammaRed        uint32
	// GammaGreen      uint32
	// GammaBlue       uint32
	// Intent          uint32
	// ProfileData     uint32
	// ProfileSize     uint32
	// Reserved        uint32
	// RGBQuads        [8]byte
}

type Bitmap struct {
	Header
	InfoHeader
	Image []byte
}

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
	var info Info
	err := binary.Read(r, binary.LittleEndian, &info)

	if err != nil {
		panic(err)
	}

	infoHeader.Info = info
	return infoHeader
}

func Rotate(bmp Bitmap) Bitmap {
	var rotated Bitmap
	rotated.Header = bmp.Header
	rotated.InfoHeader = bmp.InfoHeader
	rotated.Width, rotated.Height = rotated.Height, rotated.Width

	width := bmp.Width
	height := bmp.Height
	newWidth := rotated.Width
	newHeight := rotated.Height
	oldPad := (4 - ((width * 3) % 4)) % 4
	newPad := (4 - ((newWidth * 3) % 4)) % 4

	rotated.InfoHeader.ImageSize = uint32((3*newWidth + newPad) * newHeight)
	rotated.Image = make([]byte, rotated.ImageSize)

	for r := int32(0); r < height; r++ {
		for c := int32(0); c < width; c++ {
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
