package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func main() {

	flag.Parse()

	for _, path := range flag.Args() {

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %v", path, err)
		}

		defer r.Close()

		im, _, err := image.Decode(r)

		if err != nil {
			log.Fatalf("Failed to decode %s, %v", path, err)
		}

		// prepare BinaryBitmap
		bmp, err := gozxing.NewBinaryBitmapFromImage(im)

		if err != nil {
			log.Fatalf("Failed to create bitmap from %s, %v", path, err)
		}

		// decode image
		qrReader := qrcode.NewQRCodeReader()
		result, err := qrReader.Decode(bmp, nil)

		if err != nil {
			log.Fatalf("Failed to decode %s, %v", path, err)
		}

		fmt.Println(result)
	}
}
