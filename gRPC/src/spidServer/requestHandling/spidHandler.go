package requestHandling

import (
	"context"
	"fmt"
	"log"
	"spidServer/gps"
	pb "spidServer/requestHandling/protoBuffers"
)

func (h *Handler) GetSpidInfo(ctx context.Context, request *pb.GetSpidRequest) (*pb.GetSpidResponse, error) {
	spid, err :=  h.querySpid(request.SpidID)
	if err != nil {
		err = fmt.Errorf("failed to get spid info: %s", err)
		log.Print(err)
		return nil, err
	}
	return &pb.GetSpidResponse{
		Message: "Spid queried successfully.",
		Spid:    spid,
	}, nil
}

func (h *Handler) RegisterSpid(ctx context.Context, request *pb.RegisterSpidRequest) (*pb.RegisterSpidResponse, error) {
	spid, err :=  h.registerSpid(request.BatteryLevel, gps.FromProtoBufferEntity(request.Position))
	if err != nil {
		err = fmt.Errorf("failed to register spid: %s", err)
		log.Print(err)
		return nil, err
	}
	return &pb.RegisterSpidResponse{
		Message: "Spid registered successfully",
		Spid:    spid,
	}, nil
}

func (h *Handler) UpdateSpid(ctx context.Context, request *pb.UpdateSpidRequest) (*pb.UpdateSpidResponse, error) {
	err := h.updateSpid(request.Spid)
	if err != nil {
		err = fmt.Errorf("failed to update spid position: %s", err)
		log.Print(err)
		return nil, err
	}
	return &pb.UpdateSpidResponse{
		Message: "Spid position updated successfully.",
		Spid:    request.Spid,
	}, nil
}

func (h *Handler) DeleteSpid(ctx context.Context, request *pb.DeleteSpidRequest) (*pb.DeleteSpidResponse, error) {
	spid, err := h.deleteSpid(request.SpidID)
	if err != nil {
		err = fmt.Errorf("failed to delete spid: %s", err)
		log.Print(err)
		return nil, err
	}
	return &pb.DeleteSpidResponse{
		Message: "Spid deleted successfully.",
		Spid: spid,
	}, nil
}

func (h *Handler) AddRemoteSpid(ctx context.Context, request *pb.AddRemoteSpidRequest) (*pb.AddRemoteSpidResponse, error) {
	err := h.addRemoteSpid(request.Spid)
	if err != nil {
		err = fmt.Errorf("failed to add remote spid: %s", err)
		log.Print(err)
		return nil, err
	}
	return &pb.AddRemoteSpidResponse{
		Message: "Spid added remotely successfully.",
	}, nil
}

func (h *Handler) UpdateRemoteSpid(ctx context.Context, request *pb.UpdateRemoteSpidRequest) (*pb.UpdateRemoteSpidResponse, error) {
	err := h.updateRemoteSpid(request.Spid)
	if err != nil {
		err = fmt.Errorf("failed to update remote spid: %s", err)
		log.Print(err)
		return nil, err
	}
	return &pb.UpdateRemoteSpidResponse{
		Message: "Spid updated remotely successfully.",
	}, nil
}

func (h *Handler) RemoveRemoteSpid(ctx context.Context, request *pb.RemoveRemoteSpidRequest) (*pb.RemoveRemoteSpidResponse, error) {
	err := h.removeRemoteSpid(request.SpidID)
	if err != nil {
		err = fmt.Errorf("failed to remove remote spid: %s", err)
		log.Print(err)
		return nil, err
	}
	return &pb.RemoveRemoteSpidResponse{
		Message: "Spid removed remotely successfully.",
	}, nil
}
