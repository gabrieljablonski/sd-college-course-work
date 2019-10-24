package grpcWrapper

import (
	"context"
	"spidServer/requestHandling"
	pb "spidServer/requestHandling/grpc/userPB"
)

func (w *Wrapper) GetUserInfo(ctx context.Context, request *pb.ClientRequest) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Id: "1111", Type: requestHandling.GetUserInfo, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) RegisterUser(ctx context.Context, request *pb.ClientRequest) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Id: "1111", Type: requestHandling.RegisterUser, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) UpdateUserLocation(ctx context.Context, request *pb.ClientRequest) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Id: "1111", Type: requestHandling.UpdateUserLocation, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) DeleteUser(ctx context.Context, request *pb.ClientRequest) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Id: "1111", Type: requestHandling.DeleteUser, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) RequestAssociation(ctx context.Context, request *pb.ClientRequest) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Id: "1111", Type: requestHandling.RequestAssociation, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) RequestDissociation(ctx context.Context, request *pb.ClientRequest) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Id: "1111", Type: requestHandling.RequestDissociation, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) RequestSpidInfo(ctx context.Context, request *pb.ClientRequest) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Id: "1111", Type: requestHandling.RequestSpidInfo, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) RequestLockChange(ctx context.Context, request *pb.ClientRequest) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Id: "1111", Type: requestHandling.RequestLockChange, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}
