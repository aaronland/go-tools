package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func main() {

	flag.Parse()

	for _, path := range flag.Args() {

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %w", path, err)
		}

		defer r.Close()

		body, err := io.ReadAll(r)

		if err != nil {
			log.Fatalf("Failed to read %s, %w", path, err)
		}

		fc, err := geojson.UnmarshalFeatureCollection(body)

		if err != nil {
			log.Fatalf("Failed to unmarshal %s, %w", path, err)
		}

		count_points := len(fc.Features)
		points := make([]orb.Point, count_points)

		for idx, f := range fc.Features {

			orb_geom := f.Geometry

			pt := orb_geom.(orb.Point)
			points[idx] = pt
		}

		mp := orb.LineString(points)
		f := geojson.NewFeature(mp)
		f.Properties = map[string]interface{}{
			"foo": "bar",
		}

		enc_f, err := f.MarshalJSON()

		if err != nil {
			log.Fatalf("Failed to marshal new JSON for %s, %w", path, err)
		}

		fmt.Println(string(enc_f))
	}
}
