package dataio

import (
	_ "encoding/json"
	_ "fmt"
	_ "html/template"

	"database/sql"

	_ "github.com/gorilla/websocket"
	_ "github.com/lib/pq"

	// "io/ioutil"

	"log"
	_ "net/http"
	"os"
	_ "ted/pkg/handler" // TODO enable
	_ "time"
)

var DBConn *sql.DB // should be available globally

func ConnectToDB() {
	log.Println("DATABASE_URL ::", os.Getenv("DATABASE_URL"))
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	log.Println("DBConn != nil", DBConn != nil)
	DBConn = db
	log.Println("DBConn != nil", DBConn != nil)
	log.Println("DB connection established")
}
