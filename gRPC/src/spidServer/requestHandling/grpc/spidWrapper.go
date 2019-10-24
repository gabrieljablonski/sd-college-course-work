package grpc

import (
	"context"
	"spidServer/requestHandling"
)

func (w *Wrapper) GetSpidInfo(ctx context.Context, request *requestHandling.ClientRequest) (*requestHandling.ServerResponse, error) {
	return &requestHandling.ServerResponse{Id: "1111", Type: requestHandling.GetSpidInfo, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) RegisterSpid(ctx context.Context, request *requestHandling.ClientRequest) (*requestHandling.ServerResponse, error) {
	return &requestHandling.ServerResponse{Id: "1111", Type: requestHandling.RegisterSpid, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) UpdateBatteryInfo(ctx context.Context, request *requestHandling.ClientRequest) (*requestHandling.ServerResponse, error) {
	return &requestHandling.ServerResponse{Id: "1111", Type: requestHandling.UpdateBatteryInfo, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) UpdateSpidLocation(ctx context.Context, request *requestHandling.ClientRequest) (*requestHandling.ServerResponse, error) {
	return &requestHandling.ServerResponse{Id: "1111", Type: requestHandling.UpdateSpidLocation, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *Wrapper) DeleteSpid(ctx context.Context, request *requestHandling.ClientRequest) (*requestHandling.ServerResponse, error) {
	return &requestHandling.ServerResponse{Id: "1111", Type: requestHandling.DeleteSpid, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}
