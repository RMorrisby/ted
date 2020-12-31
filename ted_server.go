package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"

	"database/sql"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"

	// "io/ioutil"
	"encoding/csv"
	"log"
	"net/http"
	"os"
	_ "ted/pkg/handler" // TODO enable
	"ted/pkg/structs"
	"time"
)

var _ = websocket.PingMessage // debugging to silence the import-compiler

const (
	layoutDateISO = "2006-01-02"
	layoutTimeISO = "15:04:05"

	resultCSVFilename = "result.csv"

	resultsTable                  = "results"
	resultsTableColumnDefinitions = "id serial, name varchar(100), testrun varchar(32), category varchar(32), status varchar(32), endtime timestamp with time zone, message varchar(100)"
	resultsTableCreateSQL         = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", resultsTable, resultsTableColumnDefinitions)
	resultsTableInsertSQL         = fmt.Sprintf("INSERT INTO %s(name, testrun, category, status, endtime, message) VALUES", resultsTable)
)

var isLocal bool // cache the fact that we are running locally (or not)

// If "PORT" is set, we are not running locally
func IsTEDRunningLocally() bool {
	p := os.Getenv("PORT")
	if p != "" {
		return false
	}
	return true
}

// func getPort() string {
// 	p := os.Getenv("PORT")
// 	if p != "" {
// 		return p
// 	}
// 	return "8080"
// }

func getHostAndPort() string {
	// If "PORT" is set, we are running on Heroku
	// If not set, we are running locally (Win10)
	p := os.Getenv("PORT")

	// If Heroku, do not specify the hostname. Just return the : and the port
	if p != "" {
		return ":" + p
	}

	// If local (Win10), we should specify localhost as the host
	// This stops Win10 from asking about firewall permissions with each new build
	return "localhost:8080"
}

func main() {
	Startup()

	http.HandleFunc("/", IndexPage)
	http.HandleFunc("/is-alive", IsAliveHandler)
	http.HandleFunc("/result", ResultHandler)
	// Do everything else above this line

	log.Print("TED started")
	log.Fatal(http.ListenAndServe(getHostAndPort(), nil))
}

func Startup() {

	isLocal = IsTEDRunningLocally()
	log.Println("Running locally?", isLocal)
	InitResultsStore()
	existingResults := ReadResultsStore()
	CalcResultCounts(existingResults)
	log.Println("Startup() completed")
}

var SuccessCount int
var FailCount int

func CalcResultCounts(results []structs.Result) {
	for _, result := range results {
		IncrementCounts(result)
	}
}

func IndexPage(w http.ResponseWriter, r *http.Request) {

	// Without this, /somepaththatdoesntexist also resolves to / , which is strange & dumb
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()                       // find the time right now
	IndexPageVars := structs.PageVariables{ //store the date and time in a struct
		Date:         now.Format(layoutDateISO),
		Time:         now.Format(layoutTimeISO),
		Port:         getHostAndPort(),
		SuccessCount: SuccessCount,
		FailCount:    FailCount,
	}

	t, err := template.ParseFiles("index.html") // parse the html file index.html

	// if there is an error, log it
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, IndexPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps

	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// Handles the /isalive GET request path, returning a simple JSON object
func IsAliveHandler(w http.ResponseWriter, r *http.Request) {

	log.Print("Is-Alive called")

	data := "{\"is-alive\": true}"

	fmt.Fprintf(w, data)
}

// Handles the /result POST request path for receiving new test results
func ResultHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("/result called")
	switch r.Method {
	case "POST":

		// Now try to parse the POST body from JSON
		var result structs.Result
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields() // catch unwanted fields

		err := d.Decode(&result)
		if err != nil {
			// bad JSON or unrecognized json field
			log.Print("Bad JSON or unrecognized json field", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 'name' field is mandatory
		if result.Name == "" {
			http.Error(w, "Missing field 'name' from JSON object", http.StatusBadRequest)
			return
		}

		log.Println("Result received for test", result.Name)
		IncrementCounts(result)

		WriteResultToCSV(result)
	default:
		log.Println(r.Method, "/result called")
		fmt.Fprintf(w, "Only POST is supported for /result")
	}
}

func IncrementCounts(result structs.Result) {
	switch result.Status {
	case "PASSED":
		// log.Println("SuccessCount : ", SuccessCount)
		SuccessCount++
		// log.Println("SuccessCount : ", SuccessCount)
	case "FAILED":
		FailCount++
	default:
		log.Println("Result contained unrecognised status", result.Status)
	}
}

func InitResultsStore() {
	if isLocal {
		InitResultsCSV()
	} else {
		InitResultsDB()
	}
}

func InitResultsCSV() {

	needToWriteHeader := false
	if _, err := os.Stat(resultCSVFilename); os.IsNotExist(err) {
		abs, _ := filepath.Abs(resultCSVFilename)
		log.Println("Initialising results file", abs)
		needToWriteHeader = true
	}

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(resultCSVFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal("Failed to ", err)
	}

	// If the file is new/empty, write the header
	if needToWriteHeader {

		writer := csv.NewWriter(f)

		err = writer.Write(structs.ResultHeader())
		checkError("Cannot write header to file", err)
		writer.Flush()
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func InitResultsDB() {
	log.Println("Initialising results DB")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS results (tick timestamp)"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error creating database table: %q", err))
		return
	}

	// TODO
}

func WriteResultToStore(result structs.Result) {
	if isLocal {
		WriteResultToCSV(result)
	} else {
		WriteResultToDB(result)
	}
}
func WriteResultToCSV(result structs.Result) {
	log.Println("Will now write result to file :", result)
	// TODO use PSV instead of CSV
	// TODO don't write duplicates?
	f, err := os.OpenFile(resultCSVFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	resultArray := result.ToA()

	err = writer.Write(resultArray)
	checkError("Cannot write to file", err)

	log.Println("Wrote result to file")
}

func WriteResultToDB(result structs.Result) {
	log.Println("Writing result to DB")
	cmd := resultsTableInsertSQL + fmt.Sprintf("(%s %s %s %s %s %s)", result.Name, result.TestRunIdentifier, result.Category, result.Status, result.Timestamp, result.Message)
	if _, err := db.Exec(cmd); err != nil {
		log.Fatalf("Error writing result to DB: %q", err)
	}
	// TODO
}

func ReadResultsStore() (results []structs.Result) {
	if isLocal {
		results = ReadResultsCSV()
	} else {
		results = ReadResultsDB()
	}
	return
}
func ReadResultsCSV() []structs.Result {
	log.Println("Will now read results from file :", resultCSVFilename)
	f, err := os.Open(resultCSVFilename)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	checkError("Cannot read from file", err)
	size := len(lines)
	// log.Printf("Read %d results from file", size)

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

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", resultsTable))
	if err != nil {
		log.fatalf("Error reading results: %q", err)
	}

	log.Printf("Found %d results in DB")

	// TODO
	return nil
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
