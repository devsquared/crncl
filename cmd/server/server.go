package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/devsquared/crncl"
)

var (
	address string
	port    string
	config  string
)

type args struct {
	address string
	port    string
	config  string
}

func main() {
	// TODO: maybe we can do flags as final override, but this info also needs to be supported in config
	//  - to properly enable this, create a struct of expected args and pass those to run
	flag.StringVar(&address, "address", "", "the address for the server")
	flag.StringVar(&port, "port", "", "the port number for the server")
	flag.StringVar(&config, "config", "", "file path to config file")

	args := args{
		address: address,
		port:    port,
		config:  config,
	}

	if err := run(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args args, stdout io.Writer) error {
	stdout.Write([]byte("Welcome to crncl - a simple blog service.\n"))

	// First, parse the config.
	cfg, err := setConfigFromArgsAndFile(args)
	if err != nil {
		return fmt.Errorf("unable to parse args and config: %w", err)
	}
	fmt.Println(cfg)

	// TODO: now that we have config, pass needed data down to the pages for templating

	mux := http.NewServeMux()

	postReader := crncl.FileReader{
		Dir: "posts",
	}

	// TODO: let's look at moving this over to somewhere else; maybe pass a list of files? I want the code to read and be clear on the loading process
	//  maybe we could do a loadTemplates method or something?
	postTemplate := template.Must(template.ParseFiles("post.gohtml"))
	mux.HandleFunc("GET /posts/{slug}", crncl.PostHandler(postReader, postTemplate))

	indexTemplate := template.Must(template.ParseFiles("index.gohtml"))
	mux.HandleFunc("GET /", crncl.IndexHandler(postReader, indexTemplate)) //TODO: why does this need a post reader?

	contactTemplate := template.Must(template.ParseFiles("contact.gohtml"))
	mux.HandleFunc("GET /contact", crncl.ContactHandler(contactTemplate))

	aboutTemplate := template.Must(template.ParseFiles("about.gohtml"))
	mux.HandleFunc("GET /about", crncl.AboutHandler(aboutTemplate))

	addr := cfg.BaseURL + ":" + fmt.Sprint(cfg.Port)
	err = http.ListenAndServe(addr, mux)
	if err != nil {
		return err
	}

	return nil
}

func setConfigFromArgsAndFile(args args) (crncl.Config, error) {
	var (
		cfg crncl.Config
		err error
	)
	if args.config != "" { // if a config file is given as arg, use that
		cfg, err = crncl.GetConfigFromFile(args.config)
		if err != nil {
			return crncl.Config{}, fmt.Errorf("could not read config file: %w", err)
		}
	} else { // else just use the one at base path
		cfg, err = crncl.GetConfig()
		if err != nil {
			return crncl.Config{}, fmt.Errorf("could not read config file: %w", err)
		}
	}

	if args.address != "" { //override if given an address from args
		cfg.BaseURL = args.address
	}

	if args.port != "" { //override if given a port from args
		port, err := strconv.Atoi(args.port)
		if err != nil {
			return crncl.Config{}, fmt.Errorf("given port - %s - is not valid int: %w", args.port, err)
		}
		cfg.Port = port
	}

	return cfg, nil
}
