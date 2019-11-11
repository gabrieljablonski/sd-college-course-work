package requestHandling

import (
	"context"
	"fmt"
	"log"
	pb "spidServer/requestHandling/protoBuffers"
)

func (h *Handler) GetSpidInfo(ctx context.Context, request *pb.GetSpidRequest) (*pb.GetSpidResponse, error) {
	log.Printf("Received GetSpidInfo request: %s.", request)
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
	log.Printf("Received RegisterSpid request: %s.", request)
	err :=  h.registerSpid(request.Spid)
	if err != nil {
		err = fmt.Errorf("failed to register spid: %s", err)
		log.Print(err)
		return nil, err
	}
	response := &pb.RegisterSpidResponse{
		Message: "Spid registered successfully",
		Spid:    request.Spid,
	}
	log.Printf("Sending response: %s", response)
	return response, nil
}

func (h *Handler) UpdateSpid(ctx context.Context, request *pb.UpdateSpidRequest) (*pb.UpdateSpidResponse, error) {
	log.Printf("Received UpdateSpid request: %s.", request)
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
	log.Printf("Received DeleteSpid request: %s.", request)
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
	log.Printf("Received AddRemoteSpid request: %s.", request)
	err := h.addRemoteSpid(request.Spid)
	if err != nil {
		err = fmt.Errorf("failed to add remote spid: %s", err)
		log.Print(err)
		return nil, err
	}
	response := &pb.AddRemoteSpidResponse{
		Message: "Remote spid added successfully.",
	}
	log.Printf("Sending response: %s", response)
	return response, nil
}

func (h *Handler) UpdateRemoteSpid(ctx context.Context, request *pb.UpdateRemoteSpidRequest) (*pb.UpdateRemoteSpidResponse, error) {
	log.Printf("Received UpdateRemoteSpid request: %s.", request)
	err := h.updateRemoteSpid(request.Spid)
	if err != nil {
		err = fmt.Errorf("failed to update remote spid: %s", err)
		log.Print(err)
		return nil, err
	}
	response := &pb.UpdateRemoteSpidResponse{
		Message: "Remote spid updated successfully.",
	}
	log.Printf("Sending response: %s", response)
	return response, nil
}

func (h *Handler) RemoveRemoteSpid(ctx context.Context, request *pb.RemoveRemoteSpidRequest) (*pb.RemoveRemoteSpidResponse, error) {
	log.Printf("Received RemoveRemoteSpid request: %s.", request)
	err := h.removeRemoteSpid(request.SpidID)
	if err != nil {
		err = fmt.Errorf("failed to remove remote spid: %s", err)
		log.Print(err)
		return nil, err
	}
	response := &pb.RemoveRemoteSpidResponse{
		Message: "Remote spid removed successfully.",
	}
	log.Printf("Sending response: %s", response)
	return response, nil
}
