// jv will ensure that one or more URIs (or STDIN) contain valid JSON.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func IsValid(in io.Reader) error {

	var tmp interface{}

	dec := json.NewDecoder(in)
	err := dec.Decode(&tmp)

	if err != nil {
		return fmt.Errorf("Failed to decode JSON, %w", err)
	}

	return nil
}

func main() {

	strict := flag.Bool("string", false, "Exit with a fatal error if any errors are encountered.")

	flag.Parse()

	uris := flag.Args()

	if len(uris) == 0 {
		return
	}

	if uris[0] == "-" {

		err := IsValid(os.Stdin)

		if err != nil {

			log.Printf("Failed to open '%s', %v", "STDIN", err)

			if *strict {
				os.Exit(1)
			}
		}

	} else {

		for _, path := range uris {

			fh, err := os.Open(path)

			if err != nil {

				log.Printf("Failed to open '%s', %v", path, err)

				if *strict {
					os.Exit(1)
				} else {
					continue
				}
			}

			defer fh.Close()

			err = IsValid(fh)

			if err != nil {

				log.Printf("Failed to open '%s', %v", path, err)

				if *strict {
					os.Exit(1)
				} else {
					continue
				}
			}

			// log.Println(path, "OK")
		}
	}

	os.Exit(0)
}
