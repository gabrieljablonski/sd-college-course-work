package grpcWrapper

import (
	"context"
	pb "spidServer/requestHandling/grpcWrapper/spidProtoBuffers"
)

func (w *Wrapper) GetUserInfo(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return &pb.GetUserResponse{}, nil
}

func (w *Wrapper) RegisterUser(ctx context.Context, request *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	return &pb.RegisterUserResponse{}, nil
}

func (w *Wrapper) UpdateUserLocation(ctx context.Context, request *pb.UpdateUserLocationRequest) (*pb.UpdateUserLocationResponse, error) {
	return &pb.UpdateUserLocationResponse{}, nil
}

func (w *Wrapper) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	return &pb.DeleteUserResponse{}, nil
}

func (w *Wrapper) RequestAssociation(ctx context.Context, request *pb.RequestAssociationRequest) (*pb.RequestAssociationResponse, error) {
	return &pb.RequestAssociationResponse{}, nil
}

func (w *Wrapper) RequestDissociation(ctx context.Context, request *pb.RequestDissociationRequest) (*pb.RequestDissociationResponse, error) {
	return &pb.RequestDissociationResponse{}, nil
}

func (w *Wrapper) RequestSpidInfo(ctx context.Context, request *pb.RequestSpidInfoRequest) (*pb.RequestSpidInfoResponse, error) {
	return &pb.RequestSpidInfoResponse{}, nil
}

func (w *Wrapper) RequestLockChange(ctx context.Context, request *pb.RequestLockChangeRequest) (*pb.RequestLockChangeResponse, error) {
	return &pb.RequestLockChangeResponse{}, nil
}
