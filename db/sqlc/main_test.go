package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var conn *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:secret@localhost:3808/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	conn, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Couldn't open connection: ", err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}
