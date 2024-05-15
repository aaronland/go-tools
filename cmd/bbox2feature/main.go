package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sfomuseum/go-geojson-tools"
)

func main() {

	str_bbox := flag.String("bbox", "", "")
	latlon := flag.Bool("latlon", true, "")

	flag.Parse()

	f, err := geojson.BoundingBoxToFeature(*str_bbox, *latlon)

	if err != nil {
		log.Fatalf("Failed to derive feature, %w", err)
	}

	body, err := f.MarshalJSON()

	if err != nil {
		log.Fatalf("Failed to marshal feature, %w", err)
	}

	fmt.Println(string(body))
}
