package gapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/oriventi/simplebank/db/sqlc"
	"github.com/oriventi/simplebank/pb"
	"github.com/oriventi/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, hashErr := util.HashPassword(req.GetPassword())
	if hashErr != nil {
		return nil, status.Errorf(codes.Internal, "could not hash password")
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}
	user, dbErr := server.store.CreateUser(ctx, arg)
	if dbErr != nil {
		if v, ok := dbErr.(*pq.Error); ok {
			switch v.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists")
			}
		}

		return nil, status.Errorf(codes.Internal, "could not connect to database")
	}

	response := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return response, nil
}
