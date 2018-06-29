package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/kirkbyers/cyoa"
)

func main() {
	filename := flag.String("file", "gopher.json", "the JSON file with the story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)
	f, err := os.Open(*filename)
	if err != nil {
		os.Exit(1)
	}

	d := json.NewDecoder(f)
	var story cyoa.Story
	if err := d.Decode(&story); err != nil {
		os.Exit(1)
	}

	fmt.Printf("%+v\n", story)
}
