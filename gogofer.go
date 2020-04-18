package gogofer

import (
	"bufio"
	"log"
	"net"
)

const (
	MaxRequestSizeBytes = 255
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

type Request struct {
	Path string
}

type Response interface {
	Data() []byte
}

type Handler interface {
	Serve(Request) Response
}

type StaticTextHandler struct {
	Message []byte
}

func (h *StaticTextHandler) Serve(req Request) Response {
	return &StaticTextResponse{}
}

type StaticTextResponse struct {
	Message []byte
}

type errorResponse struct {
	Path []byte
}

func (e *errorResponse) Data() []byte {
	g := &GopherMapEntry{
		Type:    '3',
		Display: internalServerError,
		Path:    e.Path,
	}
	return (&GopherMap{Entries: []*GopherMapEntry{g}}).Data()
}

func (s *StaticTextResponse) Data() []byte {
	return s.Message
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

	r := bufio.NewReader(c)
	line, _, err := r.ReadLine()
	if err != nil {
		log.Println("Unable to read data from connection: ", err)
	}
	req := Request{Path: string(line)}
	res := s.Handler.Serve(req)
	if res == nil {
		log.Println("Error for request: ", string(line))
		res = &errorResponse{}
	}
	if _, err = c.Write(res.Data()); err != nil {
		log.Println("Error writing: ", err)
	}
}
