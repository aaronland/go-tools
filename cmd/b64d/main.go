package main

/*

$> go run cmd/b64e/main.go foo | go run cmd/b64d/main.go -from -
foo

*/

import (
	"encoding/base64"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	var from string
	var to string

	flag.StringVar(&from, "from", "", "The source to read data from. Valid options are none (read from args), - (read from STDIN) or the path to a file to read.")
	flag.StringVar(&to, "to", "", "The target to write data to. Value options none (write to STDOUT) or the path to a file to write to.")

	flag.Parse()

	var r io.Reader
	var wr io.WriteCloser

	switch to {
	case "":
		wr = os.Stdout
	default:

		w, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			log.Fatalf("Failed to open %s for writing, %v", to, err)
		}

		wr = w
	}

	switch from {
	case "":

		r = strings.NewReader(strings.Join(flag.Args(), " "))

	case "-":

		r = os.Stdin

	default:

		os_r, err := os.Open(from)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %v", from, err)
		}

		defer os_r.Close()
		r = os_r
	}

	b64_r := base64.NewDecoder(base64.StdEncoding, r)
	_, err := io.Copy(wr, b64_r)

	if err != nil {
		log.Fatalf("Failed to copy data, %v", err)
	}

	wr.Close()
}
