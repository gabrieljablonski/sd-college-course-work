package utils

import (
	"log"
	"os"
)

func HandleFatal(err error) bool {
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func HandleCloseFile(file *os.File, filename string) {
	if err := file.Close(); err != nil {
		log.Printf("Error closing file %s.", filename)
	}
}
