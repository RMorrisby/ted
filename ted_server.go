package main

import (
	"encoding/json"
	"fmt"

	"github.com/joho/godotenv"

	"net/http"
	"ted/pkg/constants"
	"ted/pkg/dataio"
	_ "ted/pkg/handler" // TODO enable
	"ted/pkg/help"
	"ted/pkg/pages"
	"ted/pkg/structs"
	"ted/pkg/ws"
	"time"

	"github.com/gorilla/websocket"

	log "github.com/romana/rlog"
)

var _ = websocket.PingMessage // debugging to silence the import-compiler

// func getPort() string {
// 	p := os.Getenv("PORT")
// 	if p != "" {
// 		return p
// 	}
// 	return "8080"
// }

// var templates = template.Must(template.ParseFiles("index.html", "data.html", "admin.html"))

// init() is invoked before main()
// var log *logrus.Logger

func init() {

	// Read in a custom config file for rlog
	// This file controls the logging level, etc.
	log.SetConfFile(".rlog.conf")

	// godotenv loads values from .env into the system
	// They can then be read in via os.Getenv(), e.g. os.Getenv("DATABASE_URL")
	if err := godotenv.Load(); err != nil {
		log.Error("No .env file found")
	}

}

func main() {

	// Before serving the pages
	startup()

	// Page support
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	// Pages
	http.HandleFunc("/", IndexPage)
	http.HandleFunc("/data", pages.DataPage)
	http.HandleFunc("/admin", pages.AdminPage)

	// APIs
	http.HandleFunc("/is-alive", IsAliveHandler)
	http.HandleFunc("/suite", SuiteHandler) // path to POST new suites into TED
	// http.HandleFunc("/suite/exists", SuiteExistsHandler) // path to GET new suites into TED
	// http.HandleFunc("/suites", pages.DataGetAllSuites)
	http.HandleFunc("/test", TestHandler) // path to POST new tests into TED
	// http.HandleFunc("/tests", pages.DataGetAllTests)
	// http.HandleFunc("/test/<test_name>", TestReadHandler) // path to GET a test
	http.HandleFunc("/result", ResultHandler)            // path to POST new results into TED
	http.HandleFunc("/results", pages.DataGetAllResults) // get all results for the UI

	http.HandleFunc("/admin/deleteallresults", pages.AdminDeleteAllResults)
	http.HandleFunc("/admin/deletealltests", pages.AdminDeleteAllTests)
	http.HandleFunc("/admin/deleteallsuites", pages.AdminDeleteAllSuites)
	http.HandleFunc("/admin/getresultcount", pages.AdminGetResultCount)
	http.HandleFunc("/admin/getalltestruncounts", pages.AdminGetAllTestRunCounts)
	http.HandleFunc("/admin/gettestcount", pages.AdminGetTestCount)
	http.HandleFunc("/admin/getsuitecount", pages.AdminGetSuiteCount)

	// Misc
	http.HandleFunc("/favicon.ico", pages.Favicon)

	// Do everything else above this line

	log.Info("TED started")
	startReloadServer()
	// log.Fatal(http.ListenAndServe(getHostAndPort(), nil))
}

func startup() {

	help.IsLocal = help.IsTEDRunningLocally()
	log.Debug("Running locally?", help.IsLocal)
	dataio.InitDB()
	existingResults := dataio.ReadResultStore()
	CalcResultCounts(existingResults)
	log.Debug("Startup() completed")
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
		Date:         now.Format(constants.LayoutDateISO),
		Time:         now.Format(constants.LayoutTimeISO),
		Port:         help.GetHostAndPort(),
		SuccessCount: SuccessCount,
		FailCount:    FailCount,
	}

	err := pages.Templates.ExecuteTemplate(w, "index.html", IndexPageVars)

	if err != nil {
		log.Error("template executing error: ", err)
	}
}

// IsAliveHandler handles the /isalive GET request path, returning a simple JSON object
func IsAliveHandler(w http.ResponseWriter, r *http.Request) {

	log.Debug("Is-Alive called")

	data := "{\"is-alive\": true}"

	fmt.Fprintf(w, data)
}

// ResultHandler handles the /result POST request path for receiving new test results
func ResultHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("/result called")
	switch r.Method {
	case "POST":

		// Now try to parse the POST body from JSON
		var result structs.Result
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields() // catch unwanted fields

		err := d.Decode(&result)
		if err != nil {
			// bad JSON or unrecognized json field
			log.Error("Bad JSON or unrecognized json field", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result = result.Trim()

		// 'name' field is mandatory
		if result.TestName == "" {
			http.Error(w, "Missing field 'TestName' from JSON object", http.StatusBadRequest)
			return
		}

		log.Debug("Result received for test", result.TestName)
		log.Debug(result)
		IncrementCounts(result)

		// If the test is not registered, return an error
		if !dataio.TestExists(result.TestName) {
			s := "Result referred to a test that was not registered"
			log.Error(s)
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		dataio.WriteResultToStore(result)
		w.WriteHeader(http.StatusCreated) // return a 201
	default:
		log.Println(r.Method, "/result called")
		http.Error(w, "Only POST is supported for /result", http.StatusMethodNotAllowed)
	}
}

// SuiteHandler handles the /suite POST request path for receiving new test suites
func SuiteHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("/suite called")
	switch r.Method {
	case "GET":
		log.Println(r.Method, "GET /suite called")
		log.Println(r.URL)
		log.Println(r.URL.Query())
		log.Println(r.URL.Query().Get("suite"))
		log.Println(r.URL.Query().Get("suite") != "")

		name := r.URL.Query().Get("suite")
		if name == "" {
			// A suite name must be supplied
			s := "No suite name supplied to " + r.Method + " " + r.URL.RequestURI() + "; URL must be /suite?suite=___"
			log.Error(s)
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		suite := dataio.GetSuite(name)
		if suite == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Suite '" + name + "' is not registered in TED"))
		} else {
			fmt.Fprintf(w, suite.ToJSON())
		}

	case "POST":

		// Now try to parse the POST body from JSON
		var suite structs.Suite
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields() // catch unwanted fields

		err := d.Decode(&suite)
		if err != nil {
			// bad JSON or unrecognized json field
			log.Error("Bad JSON or unrecognized json field", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// result = result.Trim()

		// 'name' field is mandatory
		if suite.Name == "" {
			http.Error(w, "Missing field 'name' from JSON object", http.StatusBadRequest)
			return
		}

		log.Println("New suite received :", suite.Name)

		dataio.WriteSuiteToDBIfNew(suite)
		w.WriteHeader(http.StatusCreated) // return a 201
	default:
		log.Debug(r.Method, "/suite called")
		http.Error(w, "Only GET and POST are supported for /suite", http.StatusMethodNotAllowed)
	}
}

// TestHandler handles the /test POST request path for receiving new tests
func TestHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, "/test called")
	switch r.Method {
	case "GET":
		log.Println(r.URL)
		log.Println(r.URL.Query())
		log.Println(r.URL.Query().Get("test"))
		log.Println(r.URL.Query().Get("test") != "")

		name := r.URL.Query().Get("test")
		if name == "" {
			// A test name must be supplied
			s := "No test name supplied to " + r.Method + " " + r.URL.RequestURI() + "; URL must be /test?test=___"
			log.Error(s)
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		test := dataio.GetTest(name)
		if test == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Test '" + name + "' is not registered in TED"))
		} else {
			fmt.Fprintf(w, test.ToJSON())
		}

	case "POST":

		// Now try to parse the POST body from JSON
		var test structs.Test
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields() // catch unwanted fields

		err := d.Decode(&test)
		if err != nil {
			// bad JSON or unrecognized json field
			log.Error("Bad JSON or unrecognized json field", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// result = result.Trim()

		// 'name' field is mandatory
		if test.Name == "" {
			http.Error(w, "Missing field 'name' from JSON object", http.StatusBadRequest)
			return
		}

		log.Println("New test received :", test.Name)

		dataio.WriteTestToDBIfNew(test)
		w.WriteHeader(http.StatusCreated) // return a 201
	default:
		log.Debug(r.Method, "/test called")
		http.Error(w, "Only POST is supported for /test", http.StatusMethodNotAllowed)
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
		log.Debug("Result contained unrecognised status", result.Status)
	}
}

// func InitResultsCSV() {

// 	needToWriteHeader := false
// 	if _, err := os.Stat(resultCSVFilename); os.IsNotExist(err) {
// 		abs, _ := filepath.Abs(resultCSVFilename)
// 		log.Println("Initialising results file", abs)
// 		needToWriteHeader = true
// 	}

// 	// If the file doesn't exist, create it, or append to the file
// 	f, err := os.OpenFile(resultCSVFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

// 	if err != nil {
// 		log.Fatal("Failed to ", err)
// 	}

// 	// If the file is new/empty, write the header
// 	if needToWriteHeader {

// 		writer := csv.NewWriter(f)

// 		err = writer.Write(structs.ResultHeader())
// 		CheckError("Cannot write header to file", err)
// 		writer.Flush()
// 	}

// 	if err := f.Close(); err != nil {
// 		log.Fatal(err)
// 	}
// }

////////////////////////

// Websockety stuff

func startReloadServer() {
	ws.WSHub = ws.NewHub()
	go ws.WSHub.Run()
	http.HandleFunc("/datareload", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(ws.WSHub, w, r)
	})

	startServer()
}

func startServer() {
	// log.Fatal(http.ListenAndServe(getHostAndPort(), nil))
	err := http.ListenAndServe(help.GetHostAndPort(), nil)
	if err != nil {
		log.Critical("Failed to start up the Reload server: ", err)
		return
	}
}

// func SendReload(result structs.Result) {
// 	log.Println("Will try to send result to WS")
// 	message := result.ToJSON()
// 	messageBytes := bytes.TrimSpace([]byte(message))
// 	ws.WSHub.Broadcast <- messageBytes

// 	log.Println("Result sent to WS: ", message)
// }
