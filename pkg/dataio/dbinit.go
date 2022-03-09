package dataio

import (
	"os"

	"ted/pkg/constants"

	log "github.com/romana/rlog"
)

/*
Contains methods responsible for initialising the DB and its tables
*/

// InitDB connects to the DB and initialises all DB tables (if they're not there)
func InitDB() {
	ConnectToDB()
	InitTableSuite()
	InitTableRegisteredTest()
	InitTableResult()
}

func InitVariables() {
	LatestTestRun = GetLatestTestRun()
	LatestSuite = GetSuiteForLatestResult()
}

// TODO remove all CSV methods
func InitResultStore() {
	// if help.IsLocal {
	// 	dataio.InitResultsCSV()
	// } else {
	// InitResultTable()
	// }
}

// Initialise the results CSV. Optionally allow the calling method to insist that the header be written
// with InitResultsCSV(true)
// Otherwise, just call this with InitResultsCSV()
// func InitResultsCSV(writeHeader ...bool) {

// 	// If a boolean has been passed to this method, then it requires this method to write the header
// 	var needToWriteHeader bool
// 	if len(writeHeader) == 0 {
// 		needToWriteHeader = false
// 	} else {
// 		needToWriteHeader = true
// 	}

// 	// If the file does not exist, then we should write the header after it is created
// 	if _, err := os.Stat(constants.ResultCSVFilename); os.IsNotExist(err) {
// 		abs, _ := filepath.Abs(constants.ResultCSVFilename)
// 		log.Println("Initialising results file", abs)
// 		needToWriteHeader = true
// 	}

// 	// If the file doesn't exist, create it, or append to the file
// 	f, err := os.OpenFile(constants.ResultCSVFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

// 	if err != nil {
// 		log.Fatal("Failed to ", err)
// 	}

// 	// If the file is new/empty, write the header
// 	if needToWriteHeader {

// 		writer := csv.NewWriter(f)

// 		err = writer.Write(structs.ResultHeader())
// 		help.CheckError("Cannot write header to file", err)
// 		writer.Flush()
// 	}

// 	if err := f.Close(); err != nil {
// 		log.Fatal(err)
// 	}
// }

func InitTableResult() {
	InitTable(constants.ResultTable, constants.ResultTableCreateSQL)
}

func InitTableRegisteredTest() {
	InitTable(constants.RegisteredTestTable, constants.RegisteredTestTableCreateSQL)
}

func InitTableSuite() {
	InitTable(constants.SuiteTable, constants.SuiteTableCreateSQL)
}

func InitTableStatus() {
	InitTable(constants.StatusTable, constants.StatusTableCreateSQL)
}

// Generic method to intialise a DB table. Needs the table name and the SQL that would create the table
func InitTable(name string, createSQL string) {

	log.Println("Initialising DB table", name)
	log.Println("DBConn != nil", DBConn != nil)

	log.Println("SQL :", createSQL)
	if _, err := DBConn.Exec(createSQL); err != nil {
		log.Critical("Error creating database table with SQL %s; error: %q", createSQL, err)
		os.Exit(1)
	}
}
