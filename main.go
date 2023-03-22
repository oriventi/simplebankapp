package main

import (
	"database/sql"
	"log"

	"github.com/oriventi/simplebank/api"
	db "github.com/oriventi/simplebank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgres://root:secret@localhost:3808/simple_bank?sslmode=disable"
	serverAddr = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddr)
	if err != nil {
		log.Fatal("cannot connect to api: ", err)
	}
}
