package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"spidServer/requestHandling"
	pb "spidServer/requestHandling/spidProtoBuffers"
)

func main() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("failed", err)
	}
	s := grpc.NewServer()
	pb.RegisterSpidHandlerServer(s, &requestHandling.Handler{})
	log.Print("serving...")
	err = s.Serve(listener)
	if err != nil {
		log.Fatal("failed", err)
	}
}