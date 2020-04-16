package gogogofer

type Request struct {
	Path string
}

type responseWriter interface {
	write([]byte) (int, error)
}

type TextDocument struct {

}

type GopherType uint8

type GopherMapEntry struct {
	Type GopherType
	Data string
}

type GopherMap struct {
	Entries []GopherMapEntry
}

type Handler interface {
	Serve(responseWriter, Request)
}