package gogofer

import (
	"io/ioutil"
	"os"
)

// Respond with the contents of a file
type FileResponse struct {
	File *os.File
}

func (res *FileResponse) Data() []byte {
	bits, err := ioutil.ReadAll(res.File)
	if err != nil {
		return (&errorResponse{}).Data()
	}
	return bits
}
