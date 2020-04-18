package gogofer

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
)

type GopherType byte

type GopherMapEntry struct {
	Type    GopherType
	Display []byte
	Path    []byte
	Host    []byte
	Port    []byte
}

type GopherMap struct {
	Entries     []*GopherMapEntry
	DefaultHost []byte
	DefaultPort []byte
}

func NewGopherMap(reader io.Reader, host string, port int) *GopherMap {
	gm := &GopherMap{
		Entries:     []*GopherMapEntry{},
		DefaultHost: []byte(host),
		DefaultPort: []byte(strconv.Itoa(port)),
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		tokens := bytes.Split(scanner.Bytes(), []byte{tab})
		if len(tokens) == 0 || len(tokens[0]) < 2 {
			continue
		}
		entry := &GopherMapEntry{}

	loop:
		for i, token := range tokens {

			switch i {
			// First byte is type, rest of first token is display
			case 0:
				entry.Type = GopherType(token[0])
				entry.Display = token[1:]

				// Second token is path
			case 1:
				entry.Path = token

				// Third token is server
			case 2:
				entry.Host = token

				// Fourth token is port
			case 3:
				entry.Port = token
			default:
				break loop
			}
		}
		gm.Entries = append(gm.Entries, entry)
	}
	return gm
}

func (gm *GopherMap) Data() []byte {
	buf := &bytes.Buffer{}
	for _, g := range gm.Entries {
		buf.WriteByte(byte(g.Type))
		buf.Write(g.Display)
		buf.WriteByte(tab)
		switch g.Type {
		default:
			buf.Write(g.Path)
			buf.WriteByte(tab)
			if g.Host == nil {
				buf.Write(gm.DefaultHost)
			} else {
				buf.Write(g.Host)
			}
			buf.WriteByte(tab)
			if g.Port == nil {
				buf.Write(gm.DefaultPort)
			} else {
				buf.Write(g.Port)
			}

		case 'i':
			buf.Write(fake)
			buf.WriteByte(tab)
			buf.Write(fake)
			buf.WriteByte(tab)
			buf.WriteByte(fakePort)
		}
		buf.Write(crlf)
	}
	return buf.Bytes()
}
