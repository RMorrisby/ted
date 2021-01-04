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

func ReadResultsStore() (results []structs.Result) {
	if help.IsLocal {
		results = ReadResultsCSV()
	} else {
		results = ReadResultsDB()
	}
	return
}

func ReadResultsCSV() []structs.Result {
	log.Println("Will now read results from file :", constants.ResultCSVFilename)
	f, err := os.Open(constants.ResultCSVFilename)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	help.CheckError("Cannot read from file", err)
	size := len(lines)
	// log.Printf("Read %d results from file", size)
	if size == 0 {
		log.Fatal("Results CSV is empty - this should always have the header row")
		return nil
	}

	records := make([]structs.Result, size-1)

	// Convert each of the lines to a Result (ignoring the header line)
	for i, line := range lines[1:] {
		result := structs.NewResult(line)
		records[i] = *result // we need the * here
	}

	// debugging
	/*
		for _, r := range records {
			log.Println(r.Status)
		}*/

	return records
}

func ReadResultsDB() []structs.Result {
	log.Println("Reading results from DB")

	rows, err := DBConn.Query(fmt.Sprintf("SELECT * FROM %s", constants.ResultsTable))
	if err != nil {
		log.Fatalf("Error reading results: %q", err)
	}

	cols, _ := rows.Columns()
	log.Printf("Found %d columns in DB", len(cols))
	// log.Printf("Found %d results in DB", resultCount)

	var results []structs.Result
	for rows.Next() {

		var r structs.Result
		err = rows.Scan(&r.Name, &r.TestRunIdentifier, &r.Category, &r.Status, &r.Timestamp, &r.Message)
		if err != nil {
			log.Fatalf("Error reading row into struct: %q", err)
		}

		results = append(results, r)
	}
	return results
}
