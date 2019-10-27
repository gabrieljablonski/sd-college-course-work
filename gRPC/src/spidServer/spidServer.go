package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"spidServer/requestHandling"
	pb "spidServer/requestHandling/protoBuffers"
)

func main() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("failed", err)
	}
	basePath, err := os.Getwd()
	handler := requestHandling.NewHandler(basePath)
	s := grpc.NewServer()
	pb.RegisterSpidHandlerServer(s, &handler)
	log.Print("serving...")
	err = s.Serve(listener)
	if err != nil {
		log.Fatal("failed", err)
	}
}