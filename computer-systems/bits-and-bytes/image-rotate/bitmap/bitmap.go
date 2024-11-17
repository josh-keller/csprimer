package bitmap

import (
	"encoding/binary"
	"io"
)

type FileHeader struct {
	Signature [2]byte
	FileSize  uint32
	Reserved  uint32
	Offset    uint32
}

type BmpInfoHeader struct {
	InfoSize uint32
	BmpInfoHeaderV5
}

type BmpInfoHeaderV3 struct {
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
}

type BmpInfoHeaderV4 struct {
	BmpInfoHeaderV3
	RedMask        uint32
	GreenMask      uint32
	BlueMask       uint32
	AlphaMask      uint32
	ColorSpaceType uint32
	ColorEndpoints CiexyzTriple
	GammaRed       uint32
	GammaGreen     uint32
	GammaBlue      uint32
}

type BmpInfoHeaderV5 struct {
	BmpInfoHeaderV4
	Intent      uint32
	ProfileData uint32
	ProfileSize uint32
	Reserved    uint32
}

type CiexyzTriple struct {
	Red   Ciexyz
	Green Ciexyz
	Blue  Ciexyz
}

type Ciexyz struct {
	X uint32
	Y uint32
	Z uint32
}

type Bitmap struct {
	FileHeader
	BmpInfoHeader
	Image []byte
	_pad  int
}

func NewFromReader(r io.Reader) Bitmap {
	var bmp Bitmap

	bmp.FileHeader = readHeader(r)
	sizeBuffer := make([]byte, 4)
	r.Read(sizeBuffer)
	ihSize := binary.LittleEndian.Uint32(sizeBuffer)
	bmp.BmpInfoHeader = readInfoHeader(r, ihSize)

	bmp.Image = readPixels(r, bmp.ImageSize)
	bmp.setPad()
	return bmp
}

func (bmp *Bitmap) setPad() {
	bmp._pad = int((4 - ((bmp.Width * 3) % 4)) % 4)
}

func readHeader(r io.Reader) FileHeader {
	h := FileHeader{}
	binary.Read(r, binary.LittleEndian, &h)
	return h
}

func readInfoHeader(r io.Reader, size uint32) BmpInfoHeader {
	infoHeader := BmpInfoHeader{InfoSize: size}
	binary.Read(r, binary.LittleEndian, &infoHeader.BmpInfoHeaderV5)
	return infoHeader
}

func readPixels(r io.Reader, size uint32) []byte {
	pixels := make([]byte, size)
	for i := uint32(0); i < size; i++ {
		binary.Read(r, binary.BigEndian, &pixels[i])
	}

	return pixels
}

func (bmp Bitmap) Write(w io.Writer) error {
	err := binary.Write(w, binary.LittleEndian, &bmp.FileHeader)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, &bmp.BmpInfoHeader)
	if err != nil {
		return err
	}

	_, err = w.Write(bmp.Image)
	if err != nil {
		return err
	}

	return nil
}

func (bmp Bitmap) indexFromCoords(r, c uint) uint {
	return r*(uint(bmp.Width)*3+uint(bmp._pad)) + 3*c
}

func Rotate(bmp Bitmap) Bitmap {
	var rotated Bitmap
	rotated.FileHeader = bmp.FileHeader
	rotated.BmpInfoHeader = bmp.BmpInfoHeader
	rotated.Width, rotated.Height = rotated.Height, rotated.Width
	rotated.setPad()

	width := bmp.Width
	height := bmp.Height
	newWidth := rotated.Width
	newHeight := rotated.Height

	rotated.BmpInfoHeader.ImageSize = uint32((3*newWidth + int32(rotated._pad)) * newHeight)
	rotated.Image = make([]byte, rotated.ImageSize)

	for row := uint(0); row < uint(height); row++ {
		for col := uint(0); col < uint(width); col++ {
			newCol := row
			newRow := uint(newHeight) - 1 - col

			oldIdx := bmp.indexFromCoords(row, col)
			newIdx := rotated.indexFromCoords(newRow, newCol)
			rotated.Image[newIdx] = bmp.Image[oldIdx]
			rotated.Image[newIdx+1] = bmp.Image[oldIdx+1]
			rotated.Image[newIdx+2] = bmp.Image[oldIdx+2]
		}
	}

	return rotated
}
