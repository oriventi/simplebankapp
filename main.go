package main

import (
	"database/sql"
	"log"

	"github.com/oriventi/simplebank/api"
	db "github.com/oriventi/simplebank/db/sqlc"
	"github.com/oriventi/simplebank/util"

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
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot connect to api: ", err)
	}
}
