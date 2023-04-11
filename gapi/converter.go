package gapi

import (
	db "github.com/oriventi/simplebank/db/sqlc"
	"github.com/oriventi/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswortChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
