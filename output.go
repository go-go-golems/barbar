package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
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
