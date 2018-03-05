package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/benspotatoes/persistagram/internal/instagram"
)

func main() {
	file := flag.String("file", "tmp/file.html", "File to parse")
	flag.Parse()

	if file == nil {
		panic("file is required")
	}

	raw, err := ioutil.ReadFile(*file)
	if err != nil {
		panic(err)
	}

	parsed, err := instagram.Parse(string(raw))
	if err != nil {
		panic(err)
	}

	if parsed == nil {
		panic("unable to parse file")
	}

	log.Printf("Parsed page for author %q\n", parsed.Author)
	for path, filename := range parsed.Links {
		log.Printf("Parsed file %q with filename %q\n", path, filename)
	}
}
