package dataio

import (
	_ "encoding/json"
	_ "html/template"
	_ "path/filepath"

	_ "database/sql"

	_ "github.com/gorilla/websocket"
	_ "github.com/lib/pq"

	// "io/ioutil"

	"log"
	_ "net/http"
	"ted/pkg/constants"
	_ "ted/pkg/handler" // TODO enable
	"ted/pkg/structs"
	_ "time"
)

func ReadResultStore() (results []structs.Result) {
	// if help.IsLocal {
	// 	results = ReadResultCSV()
	// } else {
	results = ReadResultDB()
	// }
	return
}

// func ReadResultCSV() []structs.Result {
// 	log.Println("Will now read results from file :", constants.ResultCSVFilename)
// 	f, err := os.Open(constants.ResultCSVFilename)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer f.Close()

// 	lines, err := csv.NewReader(f).ReadAll()
// 	if err != nil {
// 		panic(err)
// 	}

// 	help.CheckError("Cannot read from file", err)
// 	size := len(lines)
// 	// log.Printf("Read %d results from file", size)
// 	if size == 0 {
// 		log.Fatal("Results CSV is empty - this should always have the header row")
// 		return nil
// 	}

// 	records := make([]structs.Result, size-1)

// 	// Convert each of the lines to a Result (ignoring the header line)
// 	for i, line := range lines[1:] {
// 		result := structs.NewResult(line)
// 		records[i] = *result // we need the * here
// 	}

// 	// debugging
// 	/*
// 		for _, r := range records {
// 			log.Println(r.Status)
// 		}*/

// 	return records
// }

func ReadResultDB() []structs.Result {
	log.Println("Reading results from DB")

	sql := constants.ResultTableSelectAllSQL
	log.Println("SQL :", sql)
	rows, err := DBConn.Query(sql)
	if err != nil {
		log.Fatalf("Error reading results: %q", err)
	}

	cols, _ := rows.Columns()
	log.Printf("Found %d columns in DB", len(cols))
	// log.Printf("Found %d results in DB", resultCount)

	var results []structs.Result
	for rows.Next() {

		var r structs.Result
		// var rowID int
		err = rows.Scan(&r.SuiteName, &r.Name, &r.TestRunIdentifier, &r.Status, &r.StartTimestamp, &r.EndTimestamp, &r.RanBy, &r.Message, &r.TedStatus, &r.TedNotes)
		if err != nil {
			log.Fatalf("Error reading row into struct: %q", err)
		}

		results = append(results, r)
	}
	return results
}

func SuiteExists(name string) bool {
	log.Println("Reading suites from DB")

	sql := constants.SuiteTableSelectNameSQL
	log.Println("SQL :", sql)

	var suiteNameInDB string
	// QueryRow is supposed to return an error if there was no row
	// If there was no error, then there was a row
	err := DBConn.QueryRow(sql).Scan(&suiteNameInDB)
	if err != nil {
		log.Println("Suite", name, "was not found in the DB")
		return false
	}

	if name != suiteNameInDB {
		log.Fatalf("Suite %s was returned from the DB, when we searched for suite %s", suiteNameInDB, name)
		return false
	}
	return true
}

func TestExists(name string) bool {
	log.Println("Reading tests from DB")

	sql := constants.RegisteredTestTableSelectNameSQL
	log.Println("SQL :", sql)

	var testNameInDB string
	// QueryRow is supposed to return an error if there was no row
	// If there was no error, then there was a row
	err := DBConn.QueryRow(sql).Scan(&testNameInDB)
	if err != nil {
		log.Println("Suite", name, "was not found in the DB")
		return false
	}

	if name != testNameInDB {
		log.Fatalf("Test %s was returned from the DB, when we searched for test %s", testNameInDB, name)
		return false
	}
	return true
}
