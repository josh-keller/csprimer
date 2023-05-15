package bitmap

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These are the two types of file headers
// If it begins with 00, it's a Win 1.x Device-Dependent Bitmap
// Otherwise, if it's 0x42 0x4D, it is a later version
var (
	win1XHeader = []byte{
		0x00, 0x00, // File Type (00)
		0x05, 0x00, // Width (5)
		0x01, 0x00, // Height (5)
		0x08, 0x00, // Byte Width (8)
		0x01, 0x08, // Planes, BitsPerPixel
	}

	winBmpFileHeader = []byte{
		0x42, 0x4D, // ID Field ("BM")
		0x66, 0x00, 0x00, 0x00, // Size of the file
		0x00, 0x00, 0x00, 0x00, // Unused
		0x36, 0x00, 0x00, 0x00, // Pixel array offset
	}
)

var (
	win2Header = []byte{
		0x00, 0x00, // File Type (00)
		0x0f, 0x00, // Width (15)
		0x0f, 0x00, // Height (15)
		0x30, 0x00, // Byte Width (48)
		0x01, 0x18, //
	}
)

func TestFileType(t *testing.T) {
	testCases := []struct {
		Input        *bytes.Reader
		ExpectedType BmpFileType
	}{
		{bytes.NewReader(win1XHeader), Win1XBmpFile},
		{bytes.NewReader(winBmpFileHeader), WinBmpFile},
		{bytes.NewReader([]byte{0x01, 0x02}), UnknownFileType},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.ExpectedType, FileType(tc.Input))
	}
}
