package main

import (
	"flag"
	_ "fmt"
	"github.com/soywod/file64"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	flag.Parse()
	args := flag.Args()

	infile := args[0]
	outfile := args[1]

	fh, err := os.Open(infile)

	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(fh)

	if err != nil {
		log.Fatal(err)
	}

	body := string(bytes)

	err = file64.Decode(body, outfile)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
