package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/devsquared/crncl"
)

var (
	address string
	port    string
)

func main() {
	flag.StringVar(&address, "address", "localhost", "the address for the server")
	flag.StringVar(&port, "port", "3030", "the port number for the server")

	args := []string{address, port}

	if err := run(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	stdout.Write([]byte("Welcome to crncl - a simple blog service.\n"))
	for _, arg := range args {
		stdout.Write([]byte(fmt.Sprintf("arg: %v\n", arg)))
	}

	mux := http.NewServeMux()

	postReader := crncl.FileReader{
		Dir: "posts",
	}
	postTemplate := template.Must(template.ParseFiles("post.gohtml"))
	mux.HandleFunc("GET /posts/{slug}", crncl.PostHandler(postReader, postTemplate))

	indexTemplate := template.Must(template.ParseFiles("index.gohtml"))
	mux.HandleFunc("GET /", crncl.IndexHandler(postReader, indexTemplate))

	contactTemplate := template.Must(template.ParseFiles("contact.gohtml"))
	mux.HandleFunc("GET /contact", crncl.ContactHandler(contactTemplate))

	aboutTemplate := template.Must(template.ParseFiles("about.gohtml"))
	mux.HandleFunc("GET /about", crncl.AboutHandler(aboutTemplate))

	addr := args[0] + ":" + args[1]
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
