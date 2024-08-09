package main

import (
	"flag"
	"fmt"
	
	"github.com/mmcloughlin/geohash"
)

func main() {

	flag.Parse()

	for _, hash := range flag.Args() {
		b := geohash.BoundingBox(hash)
		fmt.Printf("%06f,%06f,%06f,%06f", b.MinLat, b.MinLng,b.MaxLat, b.MaxLng)
	}

}
