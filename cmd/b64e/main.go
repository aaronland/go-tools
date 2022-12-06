package main

import (
	"flag"
	"fmt"
	"github.com/soywod/file64"
	"log"
)

func main() {

	flag.Parse()
	args := flag.Args()

	path := args[0]

	enc, err := file64.Encode(path)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(enc)
}
