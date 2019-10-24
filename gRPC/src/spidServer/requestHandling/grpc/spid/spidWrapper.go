package spid

import (
	"context"
	"spidServer/requestHandling"
	"spidServer/requestHandling/grpc"
)

func (w *grpc.Wrapper) GetSpidInfo(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: requestHandling.GetSpidInfo, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *grpc.Wrapper) RegisterSpid(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: requestHandling.RegisterSpid, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *grpc.Wrapper) UpdateBatteryInfo(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: requestHandling.UpdateBatteryInfo, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *grpc.Wrapper) UpdateSpidLocation(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: requestHandling.UpdateSpidLocation, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *grpc.Wrapper) DeleteSpid(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: requestHandling.DeleteSpid, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}
