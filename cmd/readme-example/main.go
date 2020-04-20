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

const gmap = "0Readme\tREADME.md\t%s\t%d\r\n"

// ReadmeHandler returns a static readme file for any request. This is meant
// to serve as an example of a custom handler. A an empty selector returns a
// very basic gophermap pointing to the readme.
type ReadmeHandler struct{}

func (handler *ReadmeHandler) Serve(selector gogofer.Selector) gogofer.Response {
	if selector.Path == "" {
		return &gogofer.StaticTextResponse{
			Message: []byte(fmt.Sprintf(gmap, *host, *port)),
		}
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
