package gapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/muditshukla3/simplebank/db/sqlc"
	"github.com/muditshukla3/simplebank/pb"
	"github.com/muditshukla3/simplebank/util"
	"github.com/muditshukla3/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	voilations := validateCreateUserRequest(request)
	if voilations != nil {
		return nil, invalidArgumentError(voilations)
	}
	hashedPassword, err := util.HashPassword(request.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
	}
	arg := db.CreateUserParams{
		Username: request.GetUsername(),
		Password: hashedPassword,
		FullName: request.GetFullName(),
		Email:    request.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user %s", err)
	}

	resp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return resp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (voilations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUserName(req.GetUsername()); err != nil {
		voilations = append(voilations, fieldVoilation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		voilations = append(voilations, fieldVoilation("password", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		voilations = append(voilations, fieldVoilation("full_name", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		voilations = append(voilations, fieldVoilation("email", err))
	}

	return voilations
}
