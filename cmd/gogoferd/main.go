package main

import (
	"flag"
	"fmt"
	"log"

	"os"
	"path/filepath"

	"github.com/petderek/gogofer"
)

const (
	usage = `gogoferd is a very basic server implementing a subset of the gopher protocol.`
)

var (
	dir  = flag.String("d", "", "the directory root containing files and gophermaps")
	host = flag.String("h", "127.0.0.1", "the host to listen on")
	port = flag.Int("p", 7070, "the port to listen on")
)

func main() {
	flag.Parse()
	log.Println("Starting server")

	rDir := resolveDirectory(*dir)

	h := &GogoferdHandler{
		Host:      *host,
		Port:      *port,
		Directory: rDir,
	}
	server := &gogofer.Server{
		Addr:    fmt.Sprintf("%s:%d", *host, *port),
		Handler: h,
	}
	log.Fatal(server.ListenAndServe())
}

func resolveDirectory(dir string) string {
	if dir == "" {
		path, err := os.Getwd()
		if err != nil {
			log.Fatal("Unable to set working directory: ", err)
		}
		return path
	}

	path, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal("Unable to get absolute path: ", err)
	}
	return path
}
