// when is a command line tool for printing one or more Unix timestamps as RFC3339 strings.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {

	var layout string

	flag.StringVar(&layout, "layout", "2006-01-02T15:04:05", "The format to use for parsing a time string.")

	flag.Parse()

	switch layout {
	case "ymd":
		layout = "2006-01-02"
	default:
		//
	}

	for _, str_t := range flag.Args() {

		t, err := time.Parse(layout, str_t)

		if err != nil {
			log.Fatalf("Failed to parse '%s', %v", str_t, err)
		}

		fmt.Println(t.Unix())
	}
}
