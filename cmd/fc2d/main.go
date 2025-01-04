package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/paulmach/orb/geojson"
)

func main() {

	var path_featurecollection string
	var path_dir string

	flag.StringVar(&path_featurecollection, "featurecollection", "", "The path to a GeoJSON FeatureCollection file.")
	flag.StringVar(&path_dir, "dir", "", "The path to the directory where features will be written to.")

	flag.Parse()

	r, err := os.Open(path_featurecollection)

	if err != nil {
		log.Fatalf("Failed to open microhoods, %v", err)
	}

	defer r.Close()

	body, err := io.ReadAll(r)

	if err != nil {
		log.Fatalf("Failed read microhoods, %v", err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(body)

	if err != nil {
		log.Fatalf("Failed to unmarshale microhoods, %v", err)
	}

	for _, f := range fc.Features {

		props := f.Properties
		id := props.MustInt("wof:id", -1)

		fname := fmt.Sprintf("%d.geojson", id)
		path := filepath.Join(path_dir, fname)

		enc, err := f.MarshalJSON()

		if err != nil {
			log.Fatalf("Failed to marshal JSON for %s, %v", path, err)
		}

		wr, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			log.Fatalf("Failed to open writer for %s, %v", path, err)
		}

		_, err = wr.Write(enc)

		if err != nil {
			log.Fatalf("Failed to write %s, %v", path, err)
		}

		err = wr.Close()

		if err != nil {
			log.Fatalf("Failed to close writer for %s, %v", path, err)
		}
	}
}
