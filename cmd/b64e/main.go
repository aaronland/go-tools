package main

/*

$> go run cmd/b64e/main.go foo | go run cmd/b64d/main.go -from -
foo

*/

import (
	"bufio"
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

	enc := base64.NewEncoder(base64.StdEncoding, wr)

	switch from {
	case "":

		input := []byte(strings.Join(flag.Args(), " "))
		enc.Write(input)

	case "-":

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			enc.Write(scanner.Bytes())
		}

		err := scanner.Err()

		if err != nil {
			log.Fatalf("Failed to read input, %v", err)
		}

	default:

		for _, path := range flag.Args() {
			r, err := os.Open(path)

			if err != nil {
				log.Fatalf("Failed to open %s for reading, %v", path, err)
			}

			defer r.Close()

			_, err = io.Copy(enc, r)

			if err != nil {
				log.Fatalf("Failed to copy %s to writer, %v", path, err)
			}
		}
	}

	enc.Close()
	wr.Close()
}
