package bitmap

import (
	"encoding/binary"
	"io"
)

type BmpFileType uint8

const (
	Win1xBitmapFile BmpFileType = iota
	OS2BitmapArray
	OS2v1xBitmapFile
	Win2BitmapFile
	Win3BitmapFile
	Win4BitmapFile
	WinNTBitmapFile
	UnknownFileType
)

type Win1XHeader struct {
	FileType     [2]byte
	Width        uint16
	Height       uint16
	ByteWidth    uint16
	Planes       uint8
	BitsPerPixel uint8
}

type WinBmpFileHeader struct {
	FileType     [2]byte
	FileSize     uint32
	Reserved1    uint16
	Reserved2    uint16
	BitmapOffset uint32
}

type Win2xBitmapHeader struct {
	HeaderSize   uint32 /* Size of this header in bytes */
	Width        int16  /* Image width in pixels */
	Height       int16  /* Image height in pixels */
	Planes       uint16 /* Number of color planes */
	BitsPerPixel uint16 /* Number of bits per pixel */
}

type Win2xPaletteElement struct {
	Blue  byte /* Blue component */
	Green byte /* Green component */
	Red   byte /* Red component */
}

type Win3xPaletteElement struct {
	Blue     byte /* Blue component */
	Green    byte /* Green component */
	Red      byte /* Red component */
	Reserved byte /* Padding */
}

type Win3xBitmapHeader struct {
	HeaderSize      uint32 /* Size of this header in bytes */
	Width           int32  /* Image width in pixels */
	Height          int32  /* Image height in pixels */
	Planes          uint16 /* Number of color planes */
	BitsPerPixel    uint16 /* Number of bits per pixel */
	Compression     uint32 /* Compression methods used */
	SizeOfBitmap    uint32 /* Size of bitmap in bytes */
	HorzResolution  int32  /* Horizontal resolution in pixels per meter */
	VertResolution  int32  /* Vertical resolution in pixels per meter */
	ColorsUsed      uint32 /* Number of colors in the image */
	ColorsImportant uint32 /* Minimum number of important colors */
}

/*
The size of the color palette is calculated from the BitsPerPixel value. The
color palette has 2, 16, 256, or 0 entries for a BitsPerPixel of 1, 4, 8, and
24, respectively. The number of color palette entries is calculated as follows:

NumberOfEntries = 1 << BitsPerPixel;
To detect the presence of a color palette in a BMP file (rather than just
assuming that a color palette does exist), you can calculate the number of
bytes between the bitmap header and the bitmap data and divide this number by
the size of a single palette element. Assuming that your code is compiled using
1-byte structure element alignment, the calculation is:

NumberOfEntries = (BitmapOffset - sizeof(WINBMPFILEHEADER) -
	 sizeof(WIN2XBITMAPHEADER)) / sizeof(WIN2XPALETTEELEMENT);
If NumberOfEntries is zero, there is no palette data; otherwise, you have the
number of elements in the color palette.
*/

/* WinNTBitMap
Compression indicates the type of encoding method used to compress the bitmap
data. 0 indicates that the data is uncompressed; 1 indicates that the 8-bit RLE
algorithm was used; 2 indicates that the 4-bit RLE algorithm was used; and 3
indicates that bitfields encoding was used. If the bitmap contains 16 or 32
bits per pixel, then only a Compression value of 3 is supported and the
RedMask, GreenMask, and BlueMask fields will be present following the header in
place of a color palette. If Compression is a value other than 3, then the file
is identical to a Windows 3.x BMP file.
*/
// typedef _WinNtBitfieldsMasks
// {
// 	DWORD RedMask;         /* Mask identifying bits of red component */
// 	DWORD GreenMask;       /* Mask identifying bits of green component */
// 	DWORD BlueMask;        /* Mask identifying bits of blue component */
// } WINNTBITFIELDSMASKS;

/*
It is clear that you cannot know the internal format of a BMP file based on the
file extension alone. But, fortunately, you can use a short algorithm to
determine the internal format of BMP files.

The FileType field of the file header is where we start. If these two byte
values are 424Dh ("BM"), then you have a single-image BMP file that may have
been created under Windows or OS/2. If FileType is the value 4142h ("BA"), then
you have an OS/2 bitmap array file. Other OS/2 BMP variations have the file
extensions .ICO and .PTR.

If your file type is "BM", then you must now read the Size field of the bitmap
header to determine the version of the file. Size will be 12 for Windows 2.x
BMP and OS/2 1.x BMP, 40 for Windows 3.x and Windows NT BMP, 12 to 64 for OS/2
2.x BMP, and 108 for Windows 4.x BMP. A Windows NT BMP file will always have a
Compression value of 3; otherwise, it is read as a Windows 3.x BMP file.

Note that the only difference between Windows 2.x BMP and OS/2 1.x BMP is the
data type of the Width and Height fields. For Windows 2.x, they are signed
shorts and for OS/2 1.x, they are unsigned shorts. Windows 3.x, Windows NT, and
OS/2 2.x BMP files only vary in the number of fields in the bitmap header and
in the interpretation of the Compression field.
*/

func FileType(r io.Reader) (BmpFileType, io.Reader) {
	var rawHeader io.ReadWriter
	tr := io.TeeReader(r, rawHeader)

	var fileType uint16
	err := binary.Read(tr, binary.BigEndian, fileType)
	if err != nil {
		panic(err)
	}

	if fileType == 0x00 {

		return Win1xBitmapFile
	} else if fileType == 0x4142 {
		return OS2BitmapArray
	} else if fileType == 0x424d {
		return distinguishWinBitmapTypes(r)
	} else {
		return UnknownFileType
	}
}

func distinguishWinBitmapTypes(r io.Reader) BmpFileType {

}
