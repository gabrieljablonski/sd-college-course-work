package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"spidServer/requestHandling"
	"spidServer/requestHandling/grpc/spid"
)

func main() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("failed", err)
	}
	s := grpc.NewServer()
	spid.RegisterSpidHandlerServer(s, &requestHandling.Handler{})
	log.Print("serving...")
	err = s.Serve(listener)
	if err != nil {
		log.Fatal("failed", err)
	}
	log.Print("served...")
}