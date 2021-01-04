package dataio

import (
	_ "encoding/json"
	_ "fmt"
	_ "html/template"
	"path/filepath"

	"database/sql"

	_ "github.com/gorilla/websocket"
	_ "github.com/lib/pq"

	// "io/ioutil"
	"encoding/csv"
	"log"
	_ "net/http"
	"os"
	"ted/pkg/constants"
	_ "ted/pkg/handler" // TODO enable
	"ted/pkg/help"
	"ted/pkg/structs"
	_ "time"
)

var DBConn *sql.DB // should be available globally

func ConnectToDB() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	log.Println("DBConn != nil", DBConn != nil)
	DBConn = db
	log.Println("DBConn != nil", DBConn != nil)
	log.Println("DB connection established")
}

// Initialise the results CSV. Optionally allow the calling method to insist that the header be written
// with InitResultsCSV(true)
// Otherwise, just call this with InitResultsCSV()
func InitResultsCSV(writeHeader ...bool) {

	// If a boolean has been passed to this method, then it requires this method to write the header
	var needToWriteHeader bool
	if len(writeHeader) == 0 {
		needToWriteHeader = false
	} else {
		needToWriteHeader = true
	}

	// If the file does not exist, then we should write the header after it is created
	if _, err := os.Stat(constants.ResultCSVFilename); os.IsNotExist(err) {
		abs, _ := filepath.Abs(constants.ResultCSVFilename)
		log.Println("Initialising results file", abs)
		needToWriteHeader = true
	}

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(constants.ResultCSVFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal("Failed to ", err)
	}

	// If the file is new/empty, write the header
	if needToWriteHeader {

		writer := csv.NewWriter(f)

		err = writer.Write(structs.ResultHeader())
		help.CheckError("Cannot write header to file", err)
		writer.Flush()
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func InitResultsDB() {
	log.Println("Initialising results DB")
	log.Println("DBConn != nil", DBConn != nil)

	// dataio.DBConn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	log.Fatalf("Error opening database: %q", err)
	// }

	if _, err := DBConn.Exec(constants.ResultsTableCreateSQL); err != nil {
		log.Panicf("Error creating database table with SQL %s; error: %q", constants.ResultsTableCreateSQL, err)
		log.Fatalf("Error creating database table: %q", err)
	}

	// TODO
}
