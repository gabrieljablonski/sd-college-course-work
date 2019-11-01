package requestHandling

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"spidServer/entities"
	"spidServer/gps"
	pb "spidServer/requestHandling/protoBuffers"
)

func (h *Handler) GetUserInfo(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err :=  h.queryUser(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %s", err)
	}
	return &pb.GetUserResponse{
		Message: "User queried successfully.",
		User:    user,
	}, nil
}

func (h *Handler) RegisterUser(ctx context.Context, request *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	user, err :=  h.registerUser(request.Name, gps.FromProtoBufferEntity(request.Position))
	if err != nil {
		return nil, fmt.Errorf("failed to register user")
	}
	return &pb.RegisterUserResponse{
		Message: "User registered successfully.",
		User:    user,
	}, nil
}

func (h *Handler) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user, err :=  h.queryUser(request.User.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %s", err)
	}
	err = h.updateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user position: %s", err)
	}
	return &pb.UpdateUserResponse{
		Message: "User position updated successfully.",
		User:    user,
	}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	user, err := h.queryUser(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %s", err)
	}
	err = h.deleteUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %s", err)
	}
	return &pb.DeleteUserResponse{
		Message: "Deleted user successfully.",
		User: user,
	}, nil
}

func (h *Handler) RequestAssociation(ctx context.Context, request *pb.RequestAssociationRequest) (*pb.RequestAssociationResponse, error) {
	pbUser, err := h.queryUser(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to request association: %s", err)
	}
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return nil, fmt.Errorf("failed to request association: %s", err)
	}
	if user.CurrentSpidID != uuid.Nil {
		return nil, fmt.Errorf("failed to request association: user is already associated to spid with id `%s`",
			                    user.CurrentSpidID.String())
	}
	pbSpid, err := h.querySpid(request.SpidID)
	if err != nil {
		return nil, fmt.Errorf("failed to request association: %s", err)
	}
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return nil, fmt.Errorf("failed to request association: %s", err)
	}
	if spid.CurrentUserID != uuid.Nil {
		return nil, fmt.Errorf("failed to request association: spid is already associated to user with id `%s`",
			                    spid.CurrentUserID.String())
	}
	user.CurrentSpidID = spid.ID
	spid.CurrentUserID = user.ID
	err = h.updateUser(user.ToProtoBufferEntity())
	if err != nil {
		return nil, fmt.Errorf("failed to request association: %s", err)
	}
	err = h.updateSpid(spid.ToProtoBufferEntity())
	if err != nil {
		// if update spid failed, rollback update user
		user.CurrentSpidID = uuid.Nil
		err2 := h.updateUser(user.ToProtoBufferEntity())
		if err2 != nil {
			// if this ever happens, server will hold inconsistent data
			return nil, fmt.Errorf("failed to request association: failed to rollback `%s`, `%s`", err, err2)
		}
		return nil, fmt.Errorf("failed to request association: %s", err)
	}
	return &pb.RequestAssociationResponse{
		Message: "Association request successful.",
		User:    user.ToProtoBufferEntity(),
	}, nil
}

func (h *Handler) RequestDissociation(ctx context.Context, request *pb.RequestDissociationRequest) (*pb.RequestDissociationResponse, error) {
	pbUser, err := h.queryUser(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to request association: %s", err)
	}
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return nil, fmt.Errorf("failed to request dissociation: %s", err)
	}
	if user.CurrentSpidID == uuid.Nil {
		return nil, fmt.Errorf("failed to request association: user is not associated to any spids")
	}
	pbSpid, err := h.querySpid(user.CurrentSpidID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to request dissociation: %s", err)
	}
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return nil, fmt.Errorf("failed to request dissociation: %s", err)
	}
	if spid.CurrentUserID != user.ID {
		return nil, fmt.Errorf("failed to request dissociation: spid is not associated to user")
	}
	user.CurrentSpidID = uuid.Nil
	spid.CurrentUserID = uuid.Nil
	err = h.updateUser(user.ToProtoBufferEntity())
	if err != nil {
		return nil, fmt.Errorf("failed to request dissociation: %s", err)
	}
	err = h.updateSpid(spid.ToProtoBufferEntity())
	if err != nil {
		// if update spid failed, rollback update user
		user.CurrentSpidID = uuid.Nil
		err2 := h.updateUser(user.ToProtoBufferEntity())
		if err2 != nil {
			// if this ever happens, server will hold inconsistent data
			return nil, fmt.Errorf("failed to request dissociation: failed to rollback `%s`, `%s`", err, err2)
		}
		return nil, fmt.Errorf("failed to request dissociation: %s", err)
	}
	return &pb.RequestDissociationResponse{
		Message: "Dissociation request successful.",
		User:    user.ToProtoBufferEntity(),
	}, nil
}

func (h *Handler) RequestSpidInfo(ctx context.Context, request *pb.RequestSpidInfoRequest) (*pb.RequestSpidInfoResponse, error) {
	pbUser, err := h.queryUser(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to request spid info: %s", err)
	}
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return nil, err
	}
	pbSpid, err := h.querySpid(request.SpidID)
	if err != nil {
		return nil, fmt.Errorf("failed to request spid info: %s", err)
	}
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return nil, err
	}
	if user.CurrentSpidID != spid.ID || spid.CurrentUserID != user.ID {
		// if this happens, it means server data is inconsistent
		return nil, fmt.Errorf("failed to request spid info: user with id `%s` not associated to spid with id `%s`", user.ID, spid.ID)
	}
	return &pb.RequestSpidInfoResponse{
		Message: "Spid info request successful.",
		Spid:    spid.ToProtoBufferEntity(),
	}, nil
}

func (h *Handler) RequestLockChange(ctx context.Context, request *pb.RequestLockChangeRequest) (*pb.RequestLockChangeResponse, error) {
	pbUser, err := h.queryUser(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to request lock change: %s", err)
	}
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return nil, err
	}
	pbSpid, err := h.querySpid(request.SpidID)
	if err != nil {
		return nil, fmt.Errorf("failed to request lock change: %s", err)
	}
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return nil, err
	}
	if user.CurrentSpidID != spid.ID || spid.CurrentUserID != user.ID {
		// if this happens, it means server data is inconsistent
		return nil, fmt.Errorf("failed to request lock change: user with id `%s` not associated to spid with id `%s`", user.ID, spid.ID)
	}
	err = spid.UpdateLockState(request.LockState)
	if err != nil {
		return nil, fmt.Errorf("failed to request lock change: %s", err)
	}
	err = h.updateSpid(spid.ToProtoBufferEntity())
	if err != nil {
		return nil, fmt.Errorf("failed to request lock change: %s", err)
	}
	return &pb.RequestLockChangeResponse{
		Message: "Lock change request successful.",
		Spid:    spid.ToProtoBufferEntity(),
	}, nil
}

func (h *Handler) AddRemoteUser(ctx context.Context, request *pb.AddRemoteUserRequest) (*pb.AddRemoteUserResponse, error) {
	err := h.addRemoteUser(request.User)
	if err != nil {
		return nil, err
	}
	return &pb.AddRemoteUserResponse{
		Message: "User added remotely successfully.",
	}, nil
}

func (h *Handler) UpdateRemoteUser(ctx context.Context, request *pb.UpdateRemoteUserRequest) (*pb.UpdateRemoteUserResponse, error) {
	err := h.updateRemoteUser(request.User)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateRemoteUserResponse{
		Message: "User updated remotely successfully.",
	}, nil
}

func (h *Handler) RemoveRemoteUser(ctx context.Context, request *pb.RemoveRemoteUserRequest) (*pb.RemoveRemoteUserResponse, error) {
	err := h.removeRemoteUser(request.User)
	if err != nil {
		return nil, err
	}
	return &pb.RemoveRemoteUserResponse{
		Message: "User removed remotely successfully.",
	}, nil
}
