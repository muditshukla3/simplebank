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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	voilations := validateLoginRequest(request)
	if voilations != nil {
		return nil, invalidArgumentError(voilations)
	}
	user, err := server.store.GetUser(ctx, request.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "username not found")
		}
		return nil, status.Errorf(codes.Internal, "error reading username %s", request.GetUsername())
	}

	err = util.CheckPassword(request.GetPassword(), user.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "username and password does not match")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username, server.config.AccessTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating access token")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creatig refresh token")
	}

	mtdt := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(context.Background(), db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error saving refresh token")
	}

	response := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return response, nil
}

func validateLoginRequest(req *pb.LoginUserRequest) (voilations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		voilations = append(voilations, fieldVoilation("password", err))
	}

	if err := val.ValidateUserName(req.GetUsername()); err != nil {
		voilations = append(voilations, fieldVoilation("username", err))
	}

	return voilations
}
