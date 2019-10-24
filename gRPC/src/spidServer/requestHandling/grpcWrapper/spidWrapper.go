package grpcWrapper

import (
	"context"
	pb "spidServer/requestHandling/grpcWrapper/spidProtoBuffers"
)

func (w *Wrapper) GetSpidInfo(ctx context.Context, request *pb.GetSpidRequest) (*pb.GetSpidResponse, error) {
	return &pb.GetSpidResponse{}, nil
}

func (w *Wrapper) RegisterSpid(ctx context.Context, request *pb.RegisterSpidRequest) (*pb.RegisterSpidResponse, error) {
	return &pb.RegisterSpidResponse{}, nil
}

func (w *Wrapper) UpdateBatteryInfo(ctx context.Context, request *pb.UpdateBatteryRequest) (*pb.UpdateBatteryResponse, error) {
	return &pb.UpdateBatteryResponse{}, nil
}

func (w *Wrapper) UpdateSpidLocation(ctx context.Context, request *pb.UpdateSpidLocationRequest) (*pb.UpdateSpidLocationResponse, error) {
	return &pb.UpdateSpidLocationResponse{}, nil
}

func (w *Wrapper) DeleteSpid(ctx context.Context, request *pb.DeleteSpidRequest) (*pb.DeleteSpidResponse, error) {
	return &pb.DeleteSpidResponse{}, nil
}
