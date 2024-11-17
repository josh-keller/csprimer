package bitmap

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

var (
	squarePixels = []byte{
		0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00,
		0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff,
		0x00, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff,
		0x00, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff,
	}
	squareInfo = []byte{
		0x28, 0x00, 0x00, 0x00, // Length of info header
		0x04, 0x00, 0x00, 0x00, // Width
		0x04, 0x00, 0x00, 0x00, // Height
		0x01, 0x00, // Planes
		0x18, 0x00, // Bits per pixel
		0, 0x00, 0x00, 0x00, // Compression
		0x30, 0x00, 0x00, 0, // Size of image (48 bytes)
		0x13, 0x0B, 0x00, 0x00, // Horizontal resolution
		0x13, 0x0B, 0x00, 0x00, // Vertical resoltuion
		0, 0x00, 0x00, 0x00, // Colors in the palette
		0, 0x00, 0x00, 0x00, // Important colors
	}
	squareHeader = []byte{
		0x42, 0x4D, // ID Field ("BM")
		0x66, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, // Unused
		0x36, 0x00, 0x00, 0x00, // Pixel array offset
	}
	squareHeadersBytes = append(squareHeader, squareInfo...)
	squareBmpBytes     = append(squareHeadersBytes, squarePixels...)
	squareBmp          = NewFromReader(bytes.NewReader(squareBmpBytes))
	rectPixels         = []byte{
		0x00, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	rectInfo = []byte{
		0x28, 0x00, 0x00, 0x00, // Length of info header
		0x05, 0x00, 0x00, 0x00, // Width
		0x03, 0x00, 0x00, 0x00, // Height
		0x01, 0x00, // Planes
		0x18, 0x00, // Bits per pixel
		0x00, 0x00, 0x00, 0x00, // Compression
		0x30, 0x00, 0x00, 0x00, // Size of image (48 bytes)
		0x13, 0x0B, 0x00, 0x00, // Horizontal resolution
		0x13, 0x0B, 0x00, 0x00, // Vertical resoltuion
		0x00, 0x00, 0x00, 0x00, // Colors in the palette
		0x00, 0x00, 0x00, 0x00, // Important colors
	}
	rectHeader = []byte{
		0x42, 0x4D, // ID Field ("BM")
		0x66, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, // Unused
		0x36, 0x00, 0x00, 0x00, // Pixel array offset
	}
	rectHeadersBytes = append(rectHeader, rectInfo...)
	rectBmpBytes     = append(rectHeadersBytes, rectPixels...)
	rectBmp          = NewFromReader(bytes.NewReader(rectBmpBytes))
)

func TestSaveBitmap(t *testing.T) {
	t.Skip()
	f, err := os.Create("test_rect.bmp")
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer f.Close()

	toWrite := [][]byte{rectHeader, rectInfo, rectPixels}

	for _, slice := range toWrite {
		n, err := f.Write(slice)
		if err != nil {
			t.Fatalf("%s", err)
		}
		if n != len(slice) {
			t.Fatalf("Expected %d bytes, but got %d", len(slice), n)
		}
	}
}

func TestNewFromReader(t *testing.T) {
	if squareBmp.Width != 4 {
		t.Errorf("Expected width of %d, got %d", 4, squareBmp.Width)
	}

	if squareBmp.Height != 4 {
		t.Errorf("Expected height of %d, got %d", 4, squareBmp.Height)
	}

	if squareBmp.InfoSize != 40 {
		t.Errorf("Expected infoSize of %d, got %d", 40, squareBmp.InfoSize)
	}

	if squareBmp.Offset != 54 {
		t.Errorf("Expected offset of %d, got %d", 54, squareBmp.InfoSize)
	}
	if !reflect.DeepEqual(squareBmp.Image[0:12], []byte{0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00}) {
		t.Errorf("Rotation expected BlBlWB, got %v", squareBmp.Image[0:12])
	}
}

func TestRotateSquare(t *testing.T) {
	rotated := Rotate(squareBmp)

	if rotated.Width != 4 {
		t.Errorf("Expected width of %d, got %d", 4, rotated.Width)
	}

	if rotated.Height != 4 {
		t.Errorf("Expected height of %d, got %d", 4, rotated.Height)
	}

	if rotated.InfoSize != 40 {
		t.Errorf("Expected infoSize of %d, got %d", 40, rotated.InfoSize)
	}

	if rotated.Offset != 54 {
		t.Errorf("Expected offset of %d, got %d", 54, rotated.InfoSize)
	}

	bottomRow := rotated.Image[0:12]
	if !reflect.DeepEqual(bottomRow, []byte{0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff}) {
		t.Errorf("Rotation expected BkBkWWRR, got %v", bottomRow)
	}

	topRightCorner := rotated.Image[45:48]
	if !reflect.DeepEqual(topRightCorner, []byte{0x00, 0xff, 0x00}) {
		t.Errorf("Rotation expected G, got %v", topRightCorner)
	}
}

func TestRotateRect(t *testing.T) {
	rotated := Rotate(rectBmp)

	if rotated.Width != 3 {
		t.Errorf("Expected width of %d, got %d", 4, rotated.Width)
	}

	if rotated.Height != 5 {
		t.Errorf("Expected height of %d, got %d", 4, rotated.Height)
	}

	if rotated.InfoSize != 40 {
		t.Errorf("Expected infoSize of %d, got %d", 40, rotated.InfoSize)
	}

	if rotated.Offset != 54 {
		t.Errorf("Expected offset of %d, got %d", 54, rotated.InfoSize)
	}

	bottomRow := rotated.Image[0:9]
	expectedBottomRow := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00}
	if !reflect.DeepEqual(bottomRow, expectedBottomRow) {
		t.Errorf("Rotation expected %v, got %v", expectedBottomRow, bottomRow)
	}

	topRow := rotated.Image[48:57]
	expectedTopRow := []byte{0x00, 0x00, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff}
	if !reflect.DeepEqual(topRow, expectedTopRow) {
		t.Errorf("Rotation expected %v, got %v", expectedTopRow, topRow)
	}
}
