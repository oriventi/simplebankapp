package gapi

import (
	db "github.com/oriventi/simplebank/db/sqlc"
	"github.com/oriventi/simplebank/pb"
	"github.com/oriventi/simplebank/token"
	"github.com/oriventi/simplebank/util"
)

// server serves http requests for banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// creates a new httpServer and sets up routing
func NewServer(config util.Config, store db.Store) (*Server, error) {

	maker, makerErr := token.NewPasetoMaker(config.TokenSymmetricKey)
	if makerErr != nil {
		return nil, makerErr
	}
	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: maker,
	}

	return server, nil
}