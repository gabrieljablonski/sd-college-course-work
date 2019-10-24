package grpcWrapper

import (
	"context"
	"spidServer/requestHandling"
	pb "spidServer/requestHandling/grpc/spidPB"
)

func (w *Wrapper) GetSpidInfo(ctx context.Context, request *pb.ClientRequest) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Id: "1111", Type: requestHandling.GetSpidInfo, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) RegisterSpid(ctx context.Context, request *pb.ClientRequest) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Id: "1111", Type: requestHandling.RegisterSpid, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) UpdateBatteryInfo(ctx context.Context, request *pb.ClientRequest) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Id: "1111", Type: requestHandling.UpdateBatteryInfo, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) UpdateSpidLocation(ctx context.Context, request *pb.ClientRequest) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Id: "1111", Type: requestHandling.UpdateSpidLocation, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) DeleteSpid(ctx context.Context, request *pb.ClientRequest) (*pb.ServerResponse, error) {
	return &pb.ServerResponse{Id: "1111", Type: requestHandling.DeleteSpid, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}
