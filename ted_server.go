package main

import (
	"fmt"

	"github.com/joho/godotenv"

	"net/http"
	"ted/pkg/constants"
	"ted/pkg/dataio"
	"ted/pkg/handler"
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
	http.HandleFunc("/history", pages.HistoryPage)
	http.HandleFunc("/admin", pages.AdminPage)

	// APIs
	http.HandleFunc("/is-alive", IsAliveHandler)
	http.HandleFunc("/suite", handler.SuiteHandler) // path to POST new suites into TED
	// http.HandleFunc("/suite/exists", SuiteExistsHandler) // path to GET new suites into TED
	// http.HandleFunc("/suites", pages.DataGetAllSuites)
	http.HandleFunc("/test", handler.TestHandler) // path to POST new tests into TED
	// http.HandleFunc("/test/<test_name>", TestReadHandler) // path to GET a test
	http.HandleFunc("/result", handler.ResultHandler)    // path to POST new results into TED
	http.HandleFunc("/results", pages.DataGetAllResults) // get all results for the UI // called by data.js

	// http.HandleFunc("/history/suite", handler.GetHistoryForSuite) // path to GET suite test history // TODO needed?

	// TODO block calls to these endpoints if they are made by a browser (but not if they are made by some JS within the page)?
	http.HandleFunc("/admin/deleteallresults", pages.AdminDeleteAllResults)
	http.HandleFunc("/admin/deletealltests", pages.AdminDeleteAllTests)
	http.HandleFunc("/admin/deleteallsuites", pages.AdminDeleteAllSuites)
	http.HandleFunc("/admin/getresultcount", pages.AdminGetResultCount)
	http.HandleFunc("/admin/getalltestruncounts", pages.AdminGetAllTestRunCounts)
	http.HandleFunc("/admin/gettestcount", pages.AdminGetTestCount)
	http.HandleFunc("/admin/getsuitecount", pages.AdminGetSuiteCount)
	http.HandleFunc("/admin/suites", pages.AdminGetAllSuites)
	http.HandleFunc("/admin/tests", pages.AdminGetAllTests)

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
	dataio.InitVariables()
	log.Debug("Startup() completed")
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
		Date:          now.Format(constants.LayoutDateISO),
		Time:          now.Format(constants.LayoutTimeISO),
		Port:          help.GetHostAndPort(),
		LatestTestRun: dataio.GetLatestTestRun(),
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

//////////////////////////////////

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
