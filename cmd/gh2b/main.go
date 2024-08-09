package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mmcloughlin/geohash"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func main() {

	var mode string

	flag.StringVar(&mode, "mode", "latlon", "Valid options are: latlon, lonlat, geojson")

	flag.Parse()

	geohashes := flag.Args()

	var fc *geojson.FeatureCollection

	if mode == "geojson" && len(geohashes) > 1 {
		fc = geojson.NewFeatureCollection()
	}

	for _, hash := range geohashes {

		b := geohash.BoundingBox(hash)

		switch mode {
		case "latlon":
			fmt.Printf("%06f,%06f,%06f,%06f", b.MinLat, b.MinLng, b.MaxLat, b.MaxLng)
		case "lonlat":
			fmt.Printf("%06f,%06f,%06f,%06f", b.MinLng, b.MinLat, b.MaxLng, b.MaxLat)
		case "geojson":

			min := orb.Point([2]float64{b.MinLng, b.MinLat})
			max := orb.Point([2]float64{b.MaxLng, b.MaxLat})
			bounds := orb.Bound{Min: min, Max: max}

			f := geojson.NewFeature(bounds)
			f.Properties["geohash"] = hash

			if len(geohashes) > 1 {
				fc.Append(f)
			} else {

				enc := json.NewEncoder(os.Stdout)
				err := enc.Encode(f)

				if err != nil {
					log.Fatalf("Failed to encode feature for hash %s, %w", hash, err)
				}
			}
		}
	}

	if mode == "geojson" && len(geohashes) > 1 {

		enc := json.NewEncoder(os.Stdout)
		err := enc.Encode(fc)

		if err != nil {
			log.Fatalf("Failed to encode feature collection for hashes, %w", err)
		}

	}
}
