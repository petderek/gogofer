package gogofer

import (
	"io/ioutil"
	"os"
	"bytes"
)

// Respond with the contents of a file
type FileResponse struct {
	File *os.File
	Text bool
}

func (res *FileResponse) Data() []byte {
	bits, err := ioutil.ReadAll(res.File)
	if err != nil {
		return (&errorResponse{}).Data()
	}
	if res.Text && !bytes.HasSuffix(bits, crlf) {
		bits = append(bits, cr, lf)
	}
	return bits
}
