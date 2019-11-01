package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"spidServer/errorHandling"
)

type FileManager struct {
	BasePath string
}

func (f FileManager) GetAbsolutePath(path string) string {
	return f.BasePath + string(os.PathSeparator) + path
}

func (f FileManager) ReadFile(path string) ([]byte, error) {
	path = f.GetAbsolutePath(path)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %s", path, err)
	}
	defer errorHandling.HandleCloseFile(file, path)

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %s", path, err)
	}
	return content, nil
}

func (f FileManager) WriteToFile(path string, content []byte) error {
	path = f.GetAbsolutePath(path)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %s", path, err)
	}
	defer errorHandling.HandleCloseFile(file, path)

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %s", path, err)
	}
	log.Printf("Wrote to file: `%s`", content)
	return nil
}


