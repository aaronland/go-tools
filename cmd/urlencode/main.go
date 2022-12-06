// urlencode URL-encodes one or more command-line arguments emitting each to STDOUT. For example:
//
//	$> ./bin/urlencode 'csv://?archive=archive.tar.gz'
//	csv%3A%2F%2F%3Farchive%3Darchive.tar.gz
//
// Or:
//
//	$> echo 'csv://?archive=archive.tar.gz' | bin/urlencode -stdin
//	csv%3A%2F%2F%3Farchive%3Darchive.tar.gz
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
)

func main() {

	stdin := flag.Bool("stdin", false, "Read input from STDIN")

	unescape := flag.Bool("unescape", false, "")

	flag.Parse()

	raw := flag.Args()

	if *stdin {

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			raw = append(raw, scanner.Text())
		}

		err := scanner.Err()

		if err != nil {
			log.Fatalf("Failed to read from STDIN, %v", err)
		}

	}

	for _, str := range raw {

		if *unescape {

			enc, err := url.QueryUnescape(str)

			if err != nil {
				log.Fatalf("Failed to unescape string, %v", err)
			}

			fmt.Println(enc)
		} else {
			enc := url.QueryEscape(str)
			fmt.Println(enc)
		}

	}
}
