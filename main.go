package main

import (
	"bytes"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/qr"
	"github.com/fogleman/gg"
	"github.com/spf13/cobra"
	"image/png"
	"os"
)

const (
	width       = 512
	height      = 128
	qrCodeWidth = 200
)

var rootCmd = &cobra.Command{Use: "app"}

var cmdGenerateQR = &cobra.Command{
	Use:   "qrcode [string]",
	Short: "Generate QR code from string",
	Args:  cobra.ExactArgs(1),
	Run:   generateQR,
}

var cmdGenerateBarcode = &cobra.Command{
	Use:   "barcode [string]",
	Short: "Generate Barcode from string",
	Args:  cobra.ExactArgs(1),
	Run:   generateBarcode,
}

func main() {
	rootCmd.AddCommand(cmdGenerateQR)
	cmdGenerateQR.Flags().StringP("output", "o", "", "output file")
	rootCmd.AddCommand(cmdGenerateBarcode)
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func generateQR(cmd *cobra.Command, args []string) {
	code, err := qr.Encode(args[0], qr.H, qr.Auto)
	if err != nil {
		fmt.Println("Failed to encode string to QR code:", err)
		os.Exit(1)
	}

	code, err = barcode.Scale(code, qrCodeWidth, qrCodeWidth)

	var buffer bytes.Buffer
	err = png.Encode(&buffer, code)
	if err != nil {
		fmt.Println("Failed to encode QR code to PNG:", err)
		os.Exit(1)
	}

	if output, _ := cmd.Flags().GetString("output"); output != "" {
		err := os.WriteFile(output, buffer.Bytes(), 0644)
		cobra.CheckErr(err)
		return
	}

	outputPng(buffer.Bytes())
}

func generateBarcode(jmd *cobra.Command, args []string) {
	bc, err := code128.Encode(args[0])
	if err != nil {
		fmt.Println("Failed to encode string to Barcode:", err)
		os.Exit(1)
	}
	bc2, err := barcode.Scale(bc, width, height)
	if err != nil {
		fmt.Println("Failed to scale Barcode:", err)
		os.Exit(1)
	}
	img := gg.NewContext(width, height)
	img.DrawRectangle(0, 0, float64(width), float64(height))
	img.Fill()
	img.DrawImage(bc2, 0, 0)

	var buffer bytes.Buffer
	err = png.Encode(&buffer, img.Image())
	if err != nil {
		fmt.Println("Failed to encode Barcode to PNG:", err)
		os.Exit(1)
	}

	outputPng(buffer.Bytes())
}
