package errorHandling

import (
	"log"
	"net"
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

func HandleCloseListener(ln net.Listener) {
	if err := ln.Close(); err != nil {
		log.Print("Error closing listener.")
	}
}
