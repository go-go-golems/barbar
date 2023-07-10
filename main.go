package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"os"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/fogleman/gg"
)

const (
	width  = 512
	height = 128
)

// serializeGrCommand serializes Graphics Protocol command.
func serializeGrCommand(cmd map[string]string, payload []byte) []byte {
	res := []byte("\033_G")
	args := []string{}
	for k, v := range cmd {
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}

	res = append(res, []byte(strings.Join(args, ","))...)
	if len(payload) > 0 {
		res = append(res, ';')
		res = append(res, payload...)
	}
	res = append(res, '\033', '\\')
	return res
}

// writeChunked writes chunked Graphics Protocol command.
func writeChunked(cmd map[string]string, data []byte) {
	encoded := base64.StdEncoding.EncodeToString(data)
	const blockSize = 4096
	for len(encoded) > 0 {
		var block string
		if len(encoded) > blockSize {
			block, encoded = encoded[:blockSize], encoded[blockSize:]
		} else {
			block, encoded = encoded, ""
			cmd["m"] = "0"
		}
		payload := []byte(block)
		_, _ = os.Stdout.Write(serializeGrCommand(cmd, payload))
		cmd = make(map[string]string)
	}
}

// outputPng outputs PNG image to Kitty Terminal.
func outputPng(s []byte) {
	cmd := map[string]string{
		"a": "T",
		"f": "100",
		"m": "1",
	}
	writeChunked(cmd, s)
}

func main() {
	// Check for command line argument
	if len(os.Args) < 2 {
		panic("Usage: go run <program_name.go> (string to encode)")
	}

	// Generate a barcode
	bc, err := code128.Encode(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Scale the barcode to fit
	bc2, err := barcode.Scale(bc, width, height)
	if err != nil {
		panic(err)
	}

	// Draw the barcode to an image
	img := gg.NewContext(width, height)
	img.DrawRectangle(0, 0, float64(width), float64(height))
	img.Fill()

	img.DrawImage(bc2, 0, 0)

	// Create a buffer to hold the PNG data
	imgBuffer := new(bytes.Buffer)
	err = png.Encode(imgBuffer, img.Image())
	if err != nil {
		panic(err)
	}

	outputPng(imgBuffer.Bytes())
}
