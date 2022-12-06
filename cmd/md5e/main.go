package main

import (
	"bufio"
	"crypto/md5"
	"io"
	"log"
	"fmt"
	"os"
	"flag"
	"strings"
)

func hash(r io.Reader) error {
	
	h := md5.New()
	
	_, err := io.Copy(h, r)

	if err != nil {
		return fmt.Errorf("Failed to copy reader, %v", err)
	}

	fmt.Printf("%x", h.Sum(nil))
	return nil
}

func main() {

	mode := flag.String("mode", "", "")
	flag.Parse()

	switch *mode {
	case "file":

		for _, path := range flag.Args(){

			r, err := os.Open(path)

			if err != nil {
				log.Fatalf("Failed to open %s, %v", path, err)
			}

			defer r.Close()
			
			err = hash(r)
		
			if err != nil {
				log.Fatalf("Failed to hash %s, %v", path, err)
			}
		}
		
	case "stdin":

		r := bufio.NewReader(os.Stdin)

		err := hash(r)
		
		if err != nil {
			log.Fatalf("Failed to hash STDIN, %v", err)
		}
		
	default:

		for _, str := range flag.Args(){

			r := strings.NewReader(str)
			err := hash(r)
		
			if err != nil {
				log.Fatalf("Failed to string '%s', %v", str, err)
			}
		}
	}
}
