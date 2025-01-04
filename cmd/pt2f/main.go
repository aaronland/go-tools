package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

func main() {

	var lat float64
	var lon float64

	flag.Float64Var(&lat, "latitude", 0, "...")
	flag.Float64Var(&lon, "longitude", 0, "...")

	flag.Parse()

	pt := orb.Point([2]float64{lon, lat})
	f := geojson.NewFeature(pt)

	f.Properties = map[string]interface{}{
		"hello": "world",
	}

	enc_f, err := f.MarshalJSON()

	if err != nil {
		log.Fatalf("Failed to marshal feature, %v", err)
	}

	fmt.Println(string(enc_f))
}
