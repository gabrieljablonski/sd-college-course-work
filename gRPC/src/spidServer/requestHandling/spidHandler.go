package requestHandling

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"spidServer/entities"
	"spidServer/gps"
	pb "spidServer/requestHandling/protoBuffers"
)

func (h *Handler) querySpid(spidID string) (spid entities.Spid, err error) {
	id, err := uuid.Parse(spidID)
	if err != nil {
		return spid, fmt.Errorf("invalid user id: %s", err)
	}
	return h.DBManager.QuerySpid(id)
}

func (h *Handler) GetSpidInfo(ctx context.Context, request *pb.GetSpidRequest) (*pb.GetSpidResponse, error) {
	spid, err :=  h.querySpid(request.SpidID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spid info: %s", err)
	}
	return &pb.GetSpidResponse{
		Message: "Spid queried successfully.",
		Spid:    spid.ToProtoBufferEntity(),
	}, nil
}

func (h *Handler) RegisterSpid(ctx context.Context, request *pb.RegisterSpidRequest) (*pb.RegisterSpidResponse, error) {
	spid, err :=  h.DBManager.RegisterSpid()
	if err != nil {
		return nil, fmt.Errorf("failed to register spid")
	}
	return &pb.RegisterSpidResponse{
		Message: "Spid registered successfully",
		Spid:    spid.ToProtoBufferEntity(),
	}, nil
}

func (h *Handler) UpdateBatteryInfo(ctx context.Context, request *pb.UpdateBatteryRequest) (*pb.UpdateBatteryResponse, error) {
	spid, err := h.querySpid(request.SpidID)
	if err != nil {
		return nil, fmt.Errorf("failed to update spid battery info: %s", err)
	}
	spid.BatteryLevel = request.BatteryLevel
	err = h.DBManager.UpdateSpid(spid)
	if err != nil {
		return nil, fmt.Errorf("failed to update spid battery info: %s", err)
	}
	return &pb.UpdateBatteryResponse{
		Message: "Battery info updated successfully.",
		Spid:    spid.ToProtoBufferEntity(),
	}, nil
}

func (h *Handler) UpdateSpidLocation(ctx context.Context, request *pb.UpdateSpidLocationRequest) (*pb.UpdateSpidLocationResponse, error) {
	spid, err := h.querySpid(request.SpidID)
	if err != nil {
		return nil, fmt.Errorf("failed to update spid location: %s", err)
	}
	spid.Location = gps.FromProtoBufferEntity(request.Location)
	err = h.DBManager.UpdateSpid(spid)
	if err != nil {
		return nil, fmt.Errorf("failed to update spid location: %s", err)
	}
	return &pb.UpdateSpidLocationResponse{
		Message: "Spid location updated successfully.",
		Spid:    spid.ToProtoBufferEntity(),
	}, nil
}

func (h *Handler) DeleteSpid(ctx context.Context, request *pb.DeleteSpidRequest) (*pb.DeleteSpidResponse, error) {
	spid, err := h.querySpid(request.SpidID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete spid: %s", err)
	}
	err = h.DBManager.DeleteSpid(spid)
	if err != nil {
		return nil, fmt.Errorf("failed to delete spid: %s", err)
	}
	return &pb.DeleteSpidResponse{
		Message: "Spid deleted successfully.",
		Spid:    spid.ToProtoBufferEntity(),
	}, nil
}
