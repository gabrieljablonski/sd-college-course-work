package requestHandling

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"spidServer/entities"
	"spidServer/gps"
	pb "spidServer/requestHandling/spidProtoBuffers"
)

func (w *Wrapper) queryUser(userID string) (user entities.User, err error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return user, fmt.Errorf("invalid user id: %s", err)
	}
	return w.DBManager.QueryUser(id)
}

func (w *Wrapper) GetUserInfo(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err :=  w.queryUser(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %s", err)
	}
	return &pb.GetUserResponse{
		Message: "User queried successfully.",
		User:    user.ToProtoBufferEntity(),
	}, nil
}

func (w *Wrapper) RegisterUser(ctx context.Context, request *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	user, err :=  w.DBManager.RegisterUser(request.UserName)
	if err != nil {
		return nil, fmt.Errorf("failed to register user")
	}
	return &pb.RegisterUserResponse{
		Message: "User registered successfully.",
		User:    user.ToProtoBufferEntity(),
	}, nil
}

func (w *Wrapper) UpdateUserLocation(ctx context.Context, request *pb.UpdateUserLocationRequest) (*pb.UpdateUserLocationResponse, error) {
	if request.Location == nil {
		return nil, fmt.Errorf("missing user location")
	}
	user, err :=  w.queryUser(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user location: %s", err)
	}
	user.Location = gps.GlobalPosition{
		Latitude:  request.Location.Latitude,
		Longitude: request.Location.Longitude,
	}
	err = w.DBManager.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user location: %s", err)
	}
	return &pb.UpdateUserLocationResponse{
		Message: "User location updated successfully.",
		User:    user.ToProtoBufferEntity(),
	}, nil
}

func (w *Wrapper) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	user, err := w.queryUser(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %s", err)
	}
	err = w.DBManager.DeleteUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %s", err)
	}
	return &pb.DeleteUserResponse{
		Message: "Deleted user successfully.",
		User: user.ToProtoBufferEntity(),
	}, nil
}

func (w *Wrapper) RequestAssociation(ctx context.Context, request *pb.RequestAssociationRequest) (*pb.RequestAssociationResponse, error) {
	user, err := w.queryUser(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to request association: %s", err)
	}
	if user.CurrentSpidID != uuid.Nil {
		return nil, fmt.Errorf("failed to request association: user is already associated to spid with id `%s`",
			                    user.CurrentSpidID.String())
	}
	spid, err := w.querySpid(request.SpidID)
	if err != nil {
		return nil, fmt.Errorf("failed to request association: %s", err)
	}
	if spid.CurrentUserID != uuid.Nil {
		return nil, fmt.Errorf("failed to request association: spid is already associated to user with id `%s`",
			                    spid.CurrentUserID.String())
	}
	user.CurrentSpidID = spid.ID
	spid.CurrentUserID = user.ID
	err = w.DBManager.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to request association: %s", err)
	}
	err = w.DBManager.UpdateSpid(spid)
	if err != nil {
		// if update spid failed, rollback update user
		user.CurrentSpidID = uuid.Nil
		err2 := w.DBManager.UpdateUser(user)
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

func (w *Wrapper) RequestDissociation(ctx context.Context, request *pb.RequestDissociationRequest) (*pb.RequestDissociationResponse, error) {
	user, err := w.queryUser(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to request association: %s", err)
	}
	if user.CurrentSpidID == uuid.Nil {
		return nil, fmt.Errorf("failed to request association: user is not associated to any spids")
	}
	spid, err := w.DBManager.QuerySpid(user.CurrentSpidID)
	if err != nil {
		return nil, fmt.Errorf("failed to request dissociation: %s", err)
	}
	if spid.CurrentUserID != user.ID {
		return nil, fmt.Errorf("failed to request dissociation: spid is not associated to user")
	}
	user.CurrentSpidID = uuid.Nil
	spid.CurrentUserID = uuid.Nil
	err = w.DBManager.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to request dissociation: %s", err)
	}
	err = w.DBManager.UpdateSpid(spid)
	if err != nil {
		// if update spid failed, rollback update user
		user.CurrentSpidID = uuid.Nil
		err2 := w.DBManager.UpdateUser(user)
		if err2 != nil {
			// if this ever happens, server will hold inconsistent data
			return nil, fmt.Errorf("failed to request dissociation: failed to rollback `%s`, `%s`", err, err2)
		}
		return nil, fmt.Errorf("failed to request association: %s", err)
	}
	return &pb.RequestDissociationResponse{
		Message: "Dissociation request successful.",
		User:    user.ToProtoBufferEntity(),
	}, nil
}

func (w *Wrapper) RequestSpidInfo(ctx context.Context, request *pb.RequestSpidInfoRequest) (*pb.RequestSpidInfoResponse, error) {
	user, err := w.queryUser(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to request spid info: %s", err)
	}
	spid, err := w.querySpid(request.SpidID)
	if err != nil {
		return nil, fmt.Errorf("failed to request spid info: %s", err)
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

func (w *Wrapper) RequestLockChange(ctx context.Context, request *pb.RequestLockChangeRequest) (*pb.RequestLockChangeResponse, error) {
	user, err := w.queryUser(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to request locak change: %s", err)
	}
	spid, err := w.querySpid(request.SpidID)
	if err != nil {
		return nil, fmt.Errorf("failed to request lock change: %s", err)
	}
	if user.CurrentSpidID != spid.ID || spid.CurrentUserID != user.ID {
		// if this happens, it means server data is inconsistent
		return nil, fmt.Errorf("failed to request spid info: user with id `%s` not associated to spid with id `%s`", user.ID, spid.ID)
	}
	err = spid.UpdateLockState(request.LockState)
	if err != nil {
		return nil, fmt.Errorf("failed to request lock change: %s", err)
	}
	err = w.DBManager.UpdateSpid(spid)
	if err != nil {
		return nil, fmt.Errorf("failed to request lock change: %s", err)
	}
	return &pb.RequestLockChangeResponse{
		Message: "Lock change request successful.",
		Spid:    spid.ToProtoBufferEntity(),
	}, nil
}
