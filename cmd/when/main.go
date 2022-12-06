// when is a command line tool for printing one or more Unix timestamps as RFC3339 strings.
package main

import (
	"log"
	"flag"
	"fmt"
	"time"
	"strconv"
)

func main() {

	flag.Parse()

	for _, str_ts := range flag.Args(){

		ts, err := strconv.ParseInt(str_ts, 10, 64)

		if err != nil {
			log.Fatalf("Failed to parse '%s', %v", str_ts, err)
		}

		t := time.Unix(ts, 0)
		
		fmt.Println(t.Format(time.RFC3339))
	}
}
