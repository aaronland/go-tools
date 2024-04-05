package main

import (
	"flag"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func main() {

	var margin int
	var height int
	var width int
	var filename string

	flag.IntVar(&margin, "margin", 1, "The unit of margin to add to the QR code.")
	flag.IntVar(&height, "height", 200, "The height in pixels of the QR code.")
	flag.IntVar(&width, "width", 200, "The width in pixels of the QR code.")
	flag.StringVar(&filename, "filename", "barcode.png", "The filename of the QR code to write.")

	flag.Parse()

	body := strings.Join(flag.Args(), " ")

	if body == "" {
		log.Fatalf("Missing body")
	}

	enc := qrcode.NewQRCodeWriter()
	formatQR := gozxing.BarcodeFormat_QR_CODE

	hints := make(map[gozxing.EncodeHintType]interface{})
	hints[gozxing.EncodeHintType_MARGIN] = margin

	im, err := enc.Encode(body, formatQR, width, height, hints)

	if err != nil {
		log.Fatalf("Failed to encode QR code, %v", err)
	}
	wr, err := os.Create(filename)

	if err != nil {
		log.Fatalf("Failed to open '%s' for writing, %v", filename, err)
	}

	err = png.Encode(wr, im)

	if err != nil {
		log.Fatalf("Failed to encode '%s',  %v", filename, err)
	}

	err = wr.Close()

	if err != nil {
		log.Fatalf("Failed to close '%s' after writing, %v", filename, err)
	}
}
