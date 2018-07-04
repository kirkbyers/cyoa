package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kirkbyers/cyoa"
)

func main() {
	port := flag.Int("port", 3000, "port to start the web server on")
	filename := flag.String("file", "gopher.json", "the JSON file with the story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		os.Exit(1)
	}

	story, err := cyoa.JSONStory(f)
	if err != nil {
		os.Exit(1)
	}

	h := cyoa.NewHandler(story)
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
