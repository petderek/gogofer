package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/petderek/gogofer"
)

const (
	txtExtension = ".txt"
	mdExtension  = ".md"
	gophermap    = "gophermap"
	badDots      = ".."
)

// GogoferdHandler expects a filesystem full of .txt files. Each directory within the filesystem
// can contain a gophermap file, which will act as the landing for that directory.
type GogoferdHandler struct {
	Directory string
	Host      string
	Port      int
}

func (g *GogoferdHandler) Serve(selector gogofer.Selector) gogofer.Response {
	path, err := g.sanitize(selector)
	if err != nil {
		log.Print(err)
		return nil
	}
	if info, err := os.Stat(path); err != nil {
		log.Println("File not found: ", path)
		return nil
	} else if info.IsDir() {
		return g.handleDir(selector)
	}

	if filepath.Ext(path) != txtExtension && filepath.Ext(path) != mdExtension {
		log.Println("File served is not a txt or md file: ", path)
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		log.Println("Unable to open file: ", err)
		return nil
	}
	return &gogofer.FileResponse{File: f}
}

func (g *GogoferdHandler) sanitize(selector gogofer.Selector) (string, error) {
	if strings.Contains(selector.Path, badDots) {
		return "", fmt.Errorf("request contains mean dots and I don't want to mess with them: %s", selector.Path)
	}
	return filepath.Join(g.Directory, selector.Path), nil
}

func (g *GogoferdHandler) handleDir(selector gogofer.Selector) gogofer.Response {
	// does it have a gophermap?
	gophermapPath := filepath.Join(g.Directory, selector.Path, gophermap)
	info, err := os.Stat(gophermapPath)

	// use that gophermap if it exists
	if err == nil && !info.IsDir() {
		f, err := os.Open(gophermapPath)
		if err != nil {
			log.Println("Error opening gophermap file: ", err)
			return nil
		}
		return gogofer.NewGopherMap(f, g.Host, g.Port)
	}

	// build the gophermap if we can't find one
	contents, err := ioutil.ReadDir(filepath.Join(g.Directory, selector.Path))
	if err != nil {
		log.Println("Unable to list the directory: ", err)
		return nil
	}

	entries := make([]*gogofer.GopherMapEntry, 0)
	byteHost := []byte(g.Host)
	bytePort := []byte(strconv.Itoa(g.Port))

	for _, file := range contents {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		newSelector := filepath.Join(selector.Path, file.Name())
		entry := &gogofer.GopherMapEntry{
			Path:    []byte(newSelector),
			Display: []byte(file.Name()),
			Host:    byteHost,
			Port:    bytePort,
		}
		if file.IsDir() {
			entry.Type = '1'
		} else if filepath.Ext(newSelector) == txtExtension || filepath.Ext(newSelector) == mdExtension {
			entry.Type = '0'
		} else {
			continue
		}
		entries = append(entries, entry)
	}

	return &gogofer.GopherMap{
		Entries:     entries,
		DefaultHost: byteHost,
		DefaultPort: bytePort,
	}
}
