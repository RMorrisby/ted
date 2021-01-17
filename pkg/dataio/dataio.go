package dataio

import (
	"database/sql"

	"os"
	_ "ted/pkg/handler" // TODO enable

	_ "github.com/lib/pq" // This import is necessary - we must use it with the _

	log "github.com/romana/rlog"
)

var DBConn *sql.DB // should be available globally

func ConnectToDB() {
	log.Println("DATABASE_URL ::", os.Getenv("DATABASE_URL"))
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Criticalf("Error opening database: %q", err)
	}
	log.Debug("DBConn != nil", DBConn != nil)
	DBConn = db
	log.Debug("DBConn != nil", DBConn != nil)
	log.Debug("DB connection established")
}
