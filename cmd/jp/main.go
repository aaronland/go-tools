package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"github.com/facebookarchive/atomicfile"
	"github.com/tidwall/pretty"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func make_pretty(raw []byte, out io.Writer) error {

	var stub interface{}

	err := json.Unmarshal(raw, &stub)

	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(stub)

	if err != nil {
		log.Fatal(err)
	}

	pr := pretty.Pretty(b)
	r := bytes.NewReader(pr)

	_, err = io.Copy(out, r)

	return err
}

func main() {

	rewrite := flag.Bool("rewrite", false, "Rewrite the JSON file in place.")
	flag.Parse()

	for _, path := range flag.Args() {

		if path == "-" {

			scanner := bufio.NewScanner(os.Stdin)

			for scanner.Scan() {

				raw := scanner.Text()
				err := make_pretty([]byte(raw), os.Stdout)

				if err != nil {
					log.Fatal(err)
				}
			}

			err := scanner.Err()

			if err != nil {
				log.Fatal(err)
			}

			os.Exit(0)
		}

		abs_path, err := filepath.Abs(path)

		if err != nil {
			log.Fatal(err)
		}

		fh, err := os.Open(abs_path)

		if err != nil {
			log.Fatal(err)
		}

		raw, err := ioutil.ReadAll(fh)

		fh.Close()

		if err != nil {
			log.Fatal(err)
		}

		var wr io.Writer
		var afh *atomicfile.File

		wr = os.Stdout

		if *rewrite {

			afh, err = atomicfile.New(abs_path, 0644)

			if err != nil {
				log.Fatal(err)
			}

			defer afh.Close()
			wr = afh
		}

		err = make_pretty(raw, wr)

		if err != nil {

			if *rewrite {
				afh.Abort()
			}

			log.Fatal(err)
		}
	}

}
