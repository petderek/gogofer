package main

import (
	"flag"
	"log"
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

	log.Fatal("not implemented")
}
