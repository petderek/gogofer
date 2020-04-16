package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

const (
	usage = `gogoferd is a very basic server implementing a subset of the gopher protocol`
	MaxWaitingConnections = 1
	MaxActiveConnections = 1
	MaxRequestSize = 255
	ErrorMessage = "3Server Error"
)

var (
	dir = flag.String("d", "", "the directory root containing files and gophermaps")
	port = flag.Int("p", 7080, "the port to listen on")

	SlicedErrorMessage = []byte(ErrorMessage)
)

var (
	workChan = make(chan *net.Conn, MaxWaitingConnections)
	quit = make(chan struct{})
)

func main() {
	flag.Parse()
	log.Println("Starting server")

	log.Fatal(ListenAndServe())
}

func ListenAndServe() error {
	tcp, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		return err
	}

	startPolling()

	for {
		c, err := tcp.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
		}
		select {
		case workChan <- &c:
			log.Println("Handling connection from ", c.RemoteAddr().String())
		case <-quit:
			return fmt.Errorf("quit signal recieved")
		default:
			log.Println("Error: exceeded waiting connections")
			c.Write(SlicedErrorMessage)
			c.Close()
		}
	}
}

func startPolling() {
	for i := 0; i < MaxActiveConnections; i++ {
		go poll()
	}
}

func poll() {
	for {
		var conn net.Conn
		select {
		case <-quit:
			return
		case c := <- workChan:
			if c == nil {
				log.Println("This shouldn't happen -- conn is nil!")
				continue
			}
			conn = *c
		}
		buf := make([]byte, MaxRequestSize)
		_, err := conn.Read(buf)
		if err != nil {
			log.Print("error reading: ", err)
			conn.Close()
			continue
		}

		conn.Close()
	}
}

