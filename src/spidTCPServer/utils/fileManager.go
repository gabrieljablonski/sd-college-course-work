package utils

import (
	"io/ioutil"
	"main/errorHandling"
	"os"
)

type FileManager struct {
	BasePath string
}

func (f FileManager) ReadFile(path string) []byte {
	path = f.BasePath + path
	file, err := os.Open(path)
	errorHandling.HandleFatal(err)
	defer errorHandling.HandleCloseFile(file, path)

	content, err := ioutil.ReadAll(file)
	errorHandling.HandleFatal(err)
	return content
}

func (f FileManager) WriteToFile(path string, content []byte) {
	path = f.BasePath + path
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	errorHandling.HandleFatal(err)
	defer errorHandling.HandleCloseFile(file, path)

	_, err = file.Write(content)
	errorHandling.HandleFatal(err)
}


