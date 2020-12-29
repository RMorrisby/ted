package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/gorilla/websocket"
	// "io/ioutil"
	"encoding/csv"
	"log"
	"net/http"
	"os"
	_ "ted/pkg/handler" // TODO enable
	"ted/pkg/help"
	"time"
)

var _ = websocket.PingMessage // debugging to silence the import-compiler

const (
	layoutDateISO = "2006-01-02"
	layoutTimeISO = "15:04:05"

	resultCSVFilename = "result.csv"
)

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

	InitResultsCSV()
	existingResults := ReadResultsCSV()
	CalcResultCounts(existingResults)

	http.HandleFunc("/", IndexPage)
	http.HandleFunc("/is-alive", IsAliveHandler)
	http.HandleFunc("/result", ResultHandler)
	// Do everything else above this line

	log.Print("TED started")
	log.Fatal(http.ListenAndServe(getHostAndPort(), nil))
}

var SuccessCount int
var FailCount int

func CalcResultCounts(results []help.ResultStruct) {
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

	now := time.Now()                    // find the time right now
	IndexPageVars := help.PageVariables{ //store the date and time in a struct
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
		var result help.ResultStruct
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

		WriteToResultsCSV(result)
	default:
		log.Println(r.Method, "/result called")
		fmt.Fprintf(w, "Only POST is supported for /result")
	}
}

func IncrementCounts(result help.ResultStruct) {
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

		err = writer.Write(help.ResultHeader())
		checkError("Cannot write header to file", err)
		writer.Flush()
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

func WriteToResultsCSV(result help.ResultStruct) {
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

func ReadResultsCSV() []help.ResultStruct {
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

	records := make([]help.ResultStruct, size-1)

	// Convert each of the lines to a ResultStruct (ignoring the header line)
	for i, line := range lines[1:] {
		result := help.NewResultStruct(line)
		records[i] = *result // we need the * here
	}

	// debugging
	/*
		for _, r := range records {
			log.Println(r.Status)
		}*/

	return records
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
