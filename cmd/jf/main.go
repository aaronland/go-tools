package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/facebookarchive/atomicfile"
	"io"
	"log"
	"os"
)

func Flatten(in io.Reader, out io.Writer) error {

	var tmp interface{}

	dec := json.NewDecoder(in)
	err := dec.Decode(&tmp)

	if err != nil {
		return fmt.Errorf("Failed to decode JSON, %w", err)
	}

	enc := json.NewEncoder(out)
	err = enc.Encode(tmp)

	if err != nil {
		log.Fatalf("Failed to encode JSON, %w", err)
	}

	return nil
}

func main() {

	input := flag.String("input", "-", "The path to the JSON file to read. If '-' then will read from STDIN.")
	output := flag.String("output", "-", "The path to the JSON file to produce. If '-' then will write to STDOUT.")

	flag.Parse()

	var r io.Reader
	var wr io.Writer

	switch *input {
	case "-":
		r = os.Stdin
	default:
		fh, err := os.Open(*input)

		if err != nil {
			log.Fatalf("Failed to open '%s', %v", *input, err)
		}

		defer fh.Close()
		r = fh
	}

	switch *output {
	case "-":
		wr = os.Stdout
	default:

		fh, err := atomicfile.New(*output, 0644)

		if err != nil {
			log.Fatalf("Failed to open '%s' for writing, %v", *output, err)
		}

		defer fh.Close()
		wr = fh
	}

	err := Flatten(r, wr)

	if err != nil {
		log.Fatal(err)
	}
}
