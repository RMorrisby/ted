package main

import (
	"encoding/json"
	"fmt"
	_ "html/template"
	_ "path/filepath"

	_ "database/sql"

	"github.com/joho/godotenv"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"

	// "io/ioutil"

	"log"
	"net/http"
	"ted/pkg/constants"
	"ted/pkg/dataio"
	_ "ted/pkg/handler" // TODO enable
	"ted/pkg/help"
	"ted/pkg/pages"
	"ted/pkg/structs"
	"ted/pkg/ws"
	"time"
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
func init() {
	// godotenv loads values from .env into the system
	// They can then be read in via os.Getenv(), e.g. os.Getenv("DATABASE_URL")
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
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
	http.HandleFunc("/result", ResultHandler) // path to POST new results into TED
	http.HandleFunc("/results", pages.DataGetAllResults)
	http.HandleFunc("/admin/deleteall", pages.AdminDeleteAll)
	http.HandleFunc("/admin/getcount", pages.AdminGetCount)
	http.HandleFunc("/admin/getalltestruncounts", pages.AdminGetAllTestRunCounts)

	// Misc
	http.HandleFunc("/favicon.ico", pages.Favicon)

	// Do everything else above this line

	log.Print("TED started")
	startReloadServer()
	// log.Fatal(http.ListenAndServe(getHostAndPort(), nil))
}

func startup() {

	help.IsLocal = help.IsTEDRunningLocally()
	log.Println("Running locally?", help.IsLocal)
	dataio.InitDB()
	existingResults := dataio.ReadResultStore()
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
		Date:         now.Format(constants.LayoutDateISO),
		Time:         now.Format(constants.LayoutTimeISO),
		Port:         help.GetHostAndPort(),
		SuccessCount: SuccessCount,
		FailCount:    FailCount,
	}

	err := pages.Templates.ExecuteTemplate(w, "index.html", IndexPageVars)

	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// IsAliveHandler handles the /isalive GET request path, returning a simple JSON object
func IsAliveHandler(w http.ResponseWriter, r *http.Request) {

	log.Print("Is-Alive called")

	data := "{\"is-alive\": true}"

	fmt.Fprintf(w, data)
}

// ResultHandler handles the /result POST request path for receiving new test results
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

		result = result.Trim()

		// 'name' field is mandatory
		if result.Name == "" {
			http.Error(w, "Missing field 'name' from JSON object", http.StatusBadRequest)
			return
		}

		log.Println("Result received for test", result.Name)
		IncrementCounts(result)

		dataio.WriteResultToStore(result)
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
		log.Fatal("Failed to start up the Reload server: ", err)
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
