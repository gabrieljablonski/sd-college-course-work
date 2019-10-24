package user

import (
	"context"
	"spidServer/requestHandling"
	"spidServer/requestHandling/grpc"
)

func (w *grpc.Wrapper) GetUserInfo(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: requestHandling.GetUserInfo, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *grpc.Wrapper) RegisterUser(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: requestHandling.RegisterUser, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *grpc.Wrapper) UpdateUserLocation(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: requestHandling.UpdateUserLocation, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *grpc.Wrapper) DeleteUser(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: requestHandling.DeleteUser, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *grpc.Wrapper) RequestAssociation(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: requestHandling.RequestAssociation, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *grpc.Wrapper) RequestDissociation(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: requestHandling.RequestDissociation, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *grpc.Wrapper) RequestSpidInfo(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: requestHandling.RequestSpidInfo, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (w *grpc.Wrapper) RequestLockChange(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: requestHandling.RequestLockChange, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}
