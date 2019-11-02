package requestHandling

import (
	"context"
	"fmt"
	"log"
	"spidServer/gps"
	pb "spidServer/requestHandling/protoBuffers"
)

func (h *Handler) GetSpidInfo(ctx context.Context, request *pb.GetSpidRequest) (*pb.GetSpidResponse, error) {
	log.Print("Received GetSpidInfo request.")
	spid, err :=  h.querySpid(request.SpidID)
	if err != nil {
		err = fmt.Errorf("failed to get spid info: %s", err)
		log.Print(err)
		return nil, err
	}
	response := &pb.GetSpidResponse{
		Message: "Spid queried successfully.",
		Spid:    spid,
	}
	log.Printf("Sending response: %s", response)
	return response, nil
}

func (h *Handler) RegisterSpid(ctx context.Context, request *pb.RegisterSpidRequest) (*pb.RegisterSpidResponse, error) {
	log.Print("Received RegisterSpid request.")
	position, err := gps.FromProtoBufferEntity(request.Position)
	if err != nil {
		err = fmt.Errorf("failed to register spid: %s", err)
		log.Print(err)
		return nil, err
	}
	spid, err :=  h.registerSpid(request.BatteryLevel, position)
	if err != nil {
		err = fmt.Errorf("failed to register spid: %s", err)
		log.Print(err)
		return nil, err
	}
	response := &pb.RegisterSpidResponse{
		Message: "Spid registered successfully",
		Spid:    spid,
	}
	log.Printf("Sending response: %s", response)
	return response, nil
}

func (h *Handler) UpdateSpid(ctx context.Context, request *pb.UpdateSpidRequest) (*pb.UpdateSpidResponse, error) {
	log.Print("Received UpdateSpid request.")
	err := h.updateSpid(request.Spid)
	if err != nil {
		err = fmt.Errorf("failed to update spid position: %s", err)
		log.Print(err)
		return nil, err
	}
	response := &pb.UpdateSpidResponse{
		Message: "Spid position updated successfully.",
		Spid:    request.Spid,
	}
	log.Printf("Sending response: %s", response)
	return response, nil
}

func (h *Handler) DeleteSpid(ctx context.Context, request *pb.DeleteSpidRequest) (*pb.DeleteSpidResponse, error) {
	log.Print("Received DeleteSpid request.")
	spid, err := h.deleteSpid(request.SpidID)
	if err != nil {
		err = fmt.Errorf("failed to delete spid: %s", err)
		log.Print(err)
		return nil, err
	}
	response := &pb.DeleteSpidResponse{
		Message: "Spid deleted successfully.",
		Spid: spid,
	}
	log.Printf("Sending response: %s", response)
	return response, nil
}

func (h *Handler) AddRemoteSpid(ctx context.Context, request *pb.AddRemoteSpidRequest) (*pb.AddRemoteSpidResponse, error) {
	log.Print("Received AddRemoteSpid request.")
	err := h.addRemoteSpid(request.Spid)
	if err != nil {
		err = fmt.Errorf("failed to add remote spid: %s", err)
		log.Print(err)
		return nil, err
	}
	response := &pb.AddRemoteSpidResponse{
		Message: "Spid added remotely successfully.",
	}
	log.Printf("Sending response: %s", response)
	return response, nil
}

func (h *Handler) UpdateRemoteSpid(ctx context.Context, request *pb.UpdateRemoteSpidRequest) (*pb.UpdateRemoteSpidResponse, error) {
	log.Print("Received UpdateRemoteSpid request.")
	err := h.updateRemoteSpid(request.Spid)
	if err != nil {
		err = fmt.Errorf("failed to update remote spid: %s", err)
		log.Print(err)
		return nil, err
	}
	response := &pb.UpdateRemoteSpidResponse{
		Message: "Spid updated remotely successfully.",
	}
	log.Printf("Sending response: %s", response)
	return response, nil
}

func (h *Handler) RemoveRemoteSpid(ctx context.Context, request *pb.RemoveRemoteSpidRequest) (*pb.RemoveRemoteSpidResponse, error) {
	log.Print("Received RemoveRemoteSpid request.")
	err := h.removeRemoteSpid(request.SpidID)
	if err != nil {
		err = fmt.Errorf("failed to remove remote spid: %s", err)
		log.Print(err)
		return nil, err
	}
	response := &pb.RemoveRemoteSpidResponse{
		Message: "Spid removed remotely successfully.",
	}
	log.Printf("Sending response: %s", response)
	return response, nil
}
