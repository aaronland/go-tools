// now is a command line tool for printing the current time as a Unix timestamp.
package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {

	flag.Parse()

	now := time.Now()
	ts := now.Unix()

	fmt.Printf("%d\n", ts)
}
