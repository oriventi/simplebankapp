package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/oriventi/simplebank/api"
	db "github.com/oriventi/simplebank/db/sqlc"
	"github.com/oriventi/simplebank/gapi"
	"github.com/oriventi/simplebank/pb"
	"github.com/oriventi/simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	config, confErr := util.LoadConfig(".")
	if confErr != nil {
		log.Fatal("cannot load data: ", confErr)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	store := db.NewStore(conn)
	RunGrpcServer(config, store)
}

func RunGinServer(config util.Config, store db.Store) {
	server, serverErr := api.NewServer(config, store)
	if serverErr != nil {
		log.Fatal("cannot create server: ", serverErr)
	}

	err := server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot connect to api: ", err)
	}
}

func RunGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("could not init gRPCServer")
	}
	gRPCServer := grpc.NewServer()

	pb.RegisterSimpleBankServer(gRPCServer, server)
	reflection.Register(gRPCServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("could not init listener")
	}

	log.Printf("starting gRPC server on %s", listener.Addr().String())
	err = gRPCServer.Serve(listener)
	if err != nil {
		log.Fatal("could not start gRPC Server")
	}
}
