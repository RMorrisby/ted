package dataio

import (
	_ "encoding/json"
	"fmt"
	_ "html/template"
	_ "path/filepath"

	_ "database/sql"

	_ "github.com/gorilla/websocket"
	_ "github.com/lib/pq"

	// "io/ioutil"
	_ "encoding/csv"
	"log"
	_ "net/http"
	"os"
	"ted/pkg/constants"
	_ "ted/pkg/handler" // TODO enable
	"ted/pkg/help"
	_ "ted/pkg/structs"
	_ "time"
)

func DeleteAllResults() (success bool, err error) {
	if help.IsLocal {
		success, err = DeleteAllResultsCSV()
	} else {
		success, err = DeleteAllResultsDB()
	}
	return
}

func DeleteAllResultsCSV() (bool, error) {

	log.Println("Will now delete results from file :", constants.ResultCSVFilename)

	// f, err := os.OpenFile(constants.ResultCSVFilename, os.O_TRUNC, perm)
	err := os.Truncate(constants.ResultCSVFilename, 0)
	if err != nil {
		return false, fmt.Errorf("could not open file %q for truncation: %v", constants.ResultCSVFilename, err)
	}

	InitResultsCSV(true)

	// if err = f.Close(); err != nil {
	// 	return fmt.Errorf("could not close file handler for %q after truncation: %v", constants.ResultCSVFilename, err)
	// }
	return true, nil
}

func DeleteAllResultsDB() (bool, error) {

	log.Println("Will now delete results from DB")

	sql := fmt.Sprintf("DELETE FROM %s", constants.ResultsTable)
	log.Println("SQL :", sql)
	r, err := DBConn.Exec(sql)
	if err != nil {
		log.Fatalf("Error deleting all results: %q", err)
		return false, err
	}

	numDeleted, err := r.RowsAffected()
	if err != nil {
		log.Fatalf("Error deleting all results: %q", err)
	}

	log.Printf("Deleted %d results from the DB", numDeleted)

	return true, nil
}
