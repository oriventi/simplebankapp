package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/oriventi/simplebank/util"
)

var testQueries *Queries
var conn *sql.DB

func TestMain(m *testing.M) {

	config, confErr := util.LoadConfig("../..")
	if confErr != nil {
		log.Fatal("Could not load config")
	}

	var err error
	conn, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Couldn't open connection: ", err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}
