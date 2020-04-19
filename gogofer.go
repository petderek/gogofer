package gogofer

import (
	"bufio"
	"log"
	"net"
)

const (
	MaxSelectorSizeBytes = 255
)

const (
	tab      = '\t'
	cr       = '\r'
	lf       = '\n'
	fakePort = '0'
)

var (
	fake                = []byte("fake")
	crlf                = []byte{cr, lf}
	internalServerError = []byte("Internal Server Error")
	notFound            = []byte("Resource Not Found")
)

type Selector struct {
	Path string
}

type Response interface {
	Data() []byte
}

type Handler interface {
	Serve(Selector) Response
}

type StaticTextHandler struct {
	Message []byte
}

func (h *StaticTextHandler) Serve(req Selector) Response {
	return &StaticTextResponse{}
}

type StaticTextResponse struct {
	Message []byte
}

func (s *StaticTextResponse) Data() []byte {
	return s.Message
}

// do not attach an error message unless you intend for it to be visible to
// clients
type errorResponse struct {
	Err error
}

func (e *errorResponse) Data() []byte {
	g := &GopherMapEntry{
		Type:    '3',
		Display: internalServerError,
	}
	if e.Err != nil {
		// eerrerror
		g.Path = []byte(e.Err.Error())
	}

	entries := make([]*GopherMapEntry, 1, 1)
	entries[0] = g

	gm := &GopherMap{
		Entries: entries,
	}

	return gm.Data()
}


type Server struct {
	Addr    string
	Handler Handler // Handler to invoke, DefaultHandler if nil
}

var (
	staticContent = []byte(`Hello, world!
If you are seeing this, try running with a custom handler!`)
)

func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if s.Addr == "" {
		addr = "0.0.0.0:70"
	}

	if s.Handler == nil {
		s.Handler = &StaticTextHandler{Message: staticContent}
	}

	tcp, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		c, err := tcp.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
		}
		go s.handle(c)
	}
}

func (s *Server) handle(c net.Conn) {
	defer c.Close()

	scan := bufio.NewScanner(c)
	scan.Buffer(nil, MaxSelectorSizeBytes)
	scan.Scan()
	if scan.Err() != nil {
		log.Println("Unable to read data from connection: ", scan.Err())
	}
	req := Selector{Path: scan.Text()}
	res := s.Handler.Serve(req)
	if res == nil {
		log.Println("Error for selector: ", scan.Text())
		res = &errorResponse{}
	}
	if _, err := c.Write(res.Data()); err != nil {
		log.Println("Error writing: ", err)
	}
}
