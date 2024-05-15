// bounds will print the bounding box (minx,miny,maxx,maxy) for one or more URIs describing GeoJSON files.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func FeaturesFromURI(ctx context.Context, uri string, featurecollection bool) ([]*geojson.Feature, error) {

	fh, err := os.Open(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open URI '%s', %v", uri, err)
	}

	defer fh.Close()

	return FeaturesFromReader(ctx, fh, featurecollection)
}

func FeaturesFromReader(ctx context.Context, r io.Reader, featurecollection bool) ([]*geojson.Feature, error) {

	body, err := io.ReadAll(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to read feature, %v", err)
	}

	return FeaturesFromBytes(ctx, body, featurecollection)
}

func FeaturesFromBytes(ctx context.Context, body []byte, featurecollection bool) ([]*geojson.Feature, error) {

	var features []*geojson.Feature

	if featurecollection {

		fc, err := geojson.UnmarshalFeatureCollection(body)

		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal featurecollection, %w", err)
		}

		features = fc.Features
	} else {
		f, err := geojson.UnmarshalFeature(body)

		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal feature, %w", err)
		}

		features = []*geojson.Feature{f}
	}

	return features, nil
}

func main() {

	latlon := flag.Bool("latlon", false, "Print bounding box as miny,minx,maxy,maxx.")
	is_featurecollection := flag.Bool("featurecollection", false, "Calculate bounds for GeoJSON FeatureCollection.")

	flag.Parse()

	uris := flag.Args()

	ctx := context.Background()

	wr := os.Stdout

	for _, path := range uris {

		features, err := FeaturesFromURI(ctx, path, *is_featurecollection)

		if err != nil {
			log.Fatalf("Failed to derive features from URI, %v", err)
		}

		var bounds orb.Bound

		for i, f := range features {

			f_geom := f.Geometry
			f_bounds := f_geom.Bound()

			if i == 0 {
				bounds = f_bounds
			} else {
				bounds = bounds.Union(f_bounds)
			}
		}

		var coords []interface{}

		if *latlon {

			coords = []interface{}{
				bounds.Min.Y(),
				bounds.Min.X(),
				bounds.Max.Y(),
				bounds.Max.X(),
			}
		} else {

			coords = []interface{}{
				bounds.Min.X(),
				bounds.Min.Y(),
				bounds.Max.X(),
				bounds.Max.Y(),
			}
		}

		fmt.Fprintf(wr, "%f,%f,%f,%f\n", coords...)
	}
}
