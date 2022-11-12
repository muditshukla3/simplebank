package gapi

import (
	"context"
	"database/sql"

	db "github.com/muditshukla3/simplebank/db/sqlc"
	"github.com/muditshukla3/simplebank/pb"
	"github.com/muditshukla3/simplebank/util"
	"github.com/muditshukla3/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}
	voilations := validateUpdateUserRequest(request)
	if voilations != nil {
		return nil, invalidArgumentError(voilations)
	}

	if authPayload.Username != request.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}
	arg := db.UpdateUserParams{
		Username: request.GetUsername(),
		FullName: sql.NullString{
			String: request.GetFullName(),
			Valid:  request.FullName != nil,
		},
		Email: sql.NullString{
			String: request.GetEmail(),
			Valid:  request.Email != nil,
		},
	}

	if request.Password != nil {
		hashedPassword, err := util.HashPassword(request.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
		}
		arg.Password = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user %s", err)
	}

	resp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return resp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (voilations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUserName(req.GetUsername()); err != nil {
		voilations = append(voilations, fieldVoilation("username", err))
	}

	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			voilations = append(voilations, fieldVoilation("password", err))
		}
	}

	if req.FullName != nil {
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			voilations = append(voilations, fieldVoilation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			voilations = append(voilations, fieldVoilation("email", err))
		}
	}

	return voilations
}
