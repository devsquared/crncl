package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/devsquared/crncl"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	mux := http.NewServeMux()

	postReader := crncl.FileReader{
		Dir: "posts",
	}
	postTemplate := template.Must(template.ParseFiles("post.gohtml"))
	mux.HandleFunc("GET /posts/{slug}", crncl.PostHandler(postReader, postTemplate))

	indexTemplate := template.Must(template.ParseFiles("index.gohtml"))
	mux.HandleFunc("GET /", crncl.IndexHandler(postReader, indexTemplate))

	err := http.ListenAndServe(":3030", mux)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
