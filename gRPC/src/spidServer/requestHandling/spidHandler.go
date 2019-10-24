package requestHandling

import (
	"context"
)

func (h *Handler) GetSpidInfo(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: GetSpidInfo, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (h *Handler) RegisterSpid(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: RegisterSpid, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (h *Handler) UpdateBatteryInfo(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: UpdateBatteryInfo, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (h *Handler) UpdateSpidLocation(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: UpdateSpidLocation, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}

func (h *Handler) DeleteSpid(ctx context.Context, request *ClientRequest) (*ServerResponse, error) {
	return &ServerResponse{Id: "1111", Type: DeleteSpid, Ok: true, Body: "{\"message\": \"oi\"}",}, nil
}
