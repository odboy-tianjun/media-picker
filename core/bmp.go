package core

import (
	"encoding/binary"
	"fmt"
	"os"
)

// BMPFileHeader represents the bitmap file header structure.
type BMPFileHeader struct {
	FileType   [2]byte
	FileSize   uint32
	Reserved1  uint16
	Reserved2  uint16
	DataOffset uint32
}

// BMPInfoHeader represents the bitmap info header structure.
type BMPInfoHeader struct {
	Size            uint32
	Width           int32
	Height          int32
	Planes          uint16
	BitsPerPixel    uint16
	Compression     uint32
	SizeImage       uint32
	XPelsPerMeter   int32
	YPelsPerMeter   int32
	ColorsUsed      uint32
	ColorsImportant uint32
}

func readBmpInfo(filePath string) (error, int32, int32) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("=== Failed to open file: %v\n", err)
		return err, 0, 0
	}
	defer file.Close()

	// Read the file header.
	var fileHeader BMPFileHeader
	if err := binary.Read(file, binary.LittleEndian, &fileHeader); err != nil {
		fmt.Printf("=== Failed to Read file: %v\n", err)
		return err, 0, 0
	}

	// Check if it's a valid BMP file by verifying the file type.
	if string(fileHeader.FileType[:]) != "BM" {
		fmt.Printf("=== Not a valid BMP file: %v\n", err)
		return err, 0, 0
	}

	// Read the info header.
	var infoHeader BMPInfoHeader
	if err := binary.Read(file, binary.LittleEndian, &infoHeader); err != nil {
		fmt.Printf("=== Failed to Read file: %v\n", err)
		return err, 0, 0
	}
	return err, infoHeader.Width, infoHeader.Height
	//fmt.Printf("=== File size: %d bytes\n", fileHeader.FileSize)
	//fmt.Printf("=== Image dimensions: %dx%d\n", )
	//fmt.Printf("=== Bits per pixel: %d\n", infoHeader.BitsPerPixel)

	// At this point, you would typically read the pixel data into a slice,
	// taking into account any padding required for alignment and the specific
	// pixel format (e.g., RGB, RGBA, indexed color). This part is omitted here.

	// You may also need to skip over a possible color palette if present.
	// The size of the palette can be inferred from the number of colors used and
	// the bits per pixel value.

	// For simplicity, we'll just read the rest of the file as raw bytes.
	//data, err := ioutil.ReadAll(file)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("=== Read %d bytes of pixel data.\n", len(data))
}
