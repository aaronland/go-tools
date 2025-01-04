package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aaronland/go-tools/constants"
	"github.com/aaronland/go-tools/geojson"
)

func main() {

	var latlon bool
	var str_bbox string

	flag.StringVar(&str_bbox, "bbox", "", "")
	flag.BoolVar(&latlon, "latlon", false, "")

	flag.Parse()

	if str_bbox == constants.STDIN {

		body, err := io.ReadAll(os.Stdin)

		if err != nil {
			log.Fatal(err)
		}

		str_bbox = string(body)
	}

	log.Println(str_bbox)

	f, err := geojson.BoundingBoxToFeature(str_bbox, latlon)

	if err != nil {
		log.Fatalf("Failed to derive feature, %w", err)
	}

	body, err := f.MarshalJSON()

	if err != nil {
		log.Fatalf("Failed to marshal feature, %w", err)
	}

	fmt.Println(string(body))
}
