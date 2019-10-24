package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"spidServer/requestHandling/grpcWrapper"
	"spidServer/requestHandling/grpcWrapper/spidPB"
)

func main() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("failed", err)
	}
	s := grpc.NewServer()
	spidPB.RegisterSpidHandlerServer(s, &grpcWrapper.Wrapper{})
	log.Print("serving...")
	err = s.Serve(listener)
	if err != nil {
		log.Fatal("failed", err)
	}
	log.Print("served...")
}