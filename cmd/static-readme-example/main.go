package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/petderek/gogofer"
)

var (
	port = flag.Int("p", 7070, "the port to listen on")
	host = flag.String("h", "127.0.0.1", "the host to listen on")
	file = flag.String("f", "README.md", "the file to serve")
)

// ReadmeHandler returns a static readme file for any request. This is meant to serve as an example of a custom handler.
type ReadmeHandler struct{}

func (handler *ReadmeHandler) Serve(_ gogofer.Selector) gogofer.Response {
	if file == nil {
		return nil
	}
	f, err := os.Open(*file)
	if err != nil {
		log.Println("Error: ", err)
		return nil
	}
	return &gogofer.FileResponse{File: f}
}

func main() {
	flag.Parse()
	server := gogofer.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: &ReadmeHandler{},
	}
	log.Fatal(server.ListenAndServe())
}
