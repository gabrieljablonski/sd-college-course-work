package utils

import (
	"io/ioutil"
	"os"
)

type FileManager struct {
	BasePath string
}

func (f FileManager) readFile(path string) []byte {
	path = f.BasePath + path
	file, err := os.Open(path)
	HandleFatal(err)
	defer HandleCloseFile(file, path)

	content, err := ioutil.ReadAll(file)
	HandleFatal(err)
	return content
}

func (f FileManager) writeToFile(path string, content []byte) {
	path = f.BasePath + path
	file, err := os.Open(path)
	HandleFatal(err)
	defer HandleCloseFile(file, path)

	_, err = file.Write(content)
	HandleFatal(err)
}


