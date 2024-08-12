package main

import (
	"flag"
	"fmt"

	"github.com/mmcloughlin/geohash"
)

func main() {

	var lat float64
	var lon float64
	var precision uint

	flag.Float64Var(&lat, "latitude", 0, "...")
	flag.Float64Var(&lon, "longitude", 0, "...")
	flag.UintVar(&precision, "precision", 5, "...")

	flag.Parse()

	fmt.Println(geohash.EncodeWithPrecision(lat, lon, precision))

}
