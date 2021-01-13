package pages

import (
	"bytes"
	_ "database/sql"
	"encoding/json"
	_ "fmt"
	_ "html/template"
	"log"
	"net/http"
	_ "os"
	_ "path/filepath"
	"strconv"
	"ted/pkg/constants"
	"ted/pkg/dataio"
	_ "ted/pkg/handler" // TODO enable
	"ted/pkg/help"
	"ted/pkg/structs"
	_ "ted/pkg/ws"
	"time"

	_ "github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

// The Admin page
func AdminPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()                       // find the time right now
	AdminPageVars := structs.PageVariables{ //store the date and time in a struct
		Date: now.Format(constants.LayoutDateISO),
		Time: now.Format(constants.LayoutTimeISO),
	}

	// ws.ServeWs(ws.WSHub, w, r)

	err := Templates.ExecuteTemplate(w, "admin.html", AdminPageVars) //execute the template and pass it the struct to fill in the gaps

	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// REST endpoint to trigger the deletion of all results
func AdminDeleteAll(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	success, err := dataio.DeleteAllResults()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	} else if success {
		w.Write([]byte(strconv.Itoa(0)))
	} else {
		http.Error(w, "ERROR - not successful but no error returned!", 500)
		return
	}
}

// REST endpoint to get the total number of results
func AdminGetCount(w http.ResponseWriter, r *http.Request) {

	log.Print("AdminGetCount called")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count := len(dataio.ReadResultsStore())
	log.Print("Total result count : ", count)
	w.Write([]byte(strconv.Itoa(count)))
}

// REST endpoint to get the total number of known test runs, and the number of results in each run
func AdminGetAllTestRunCounts(w http.ResponseWriter, r *http.Request) {

	log.Print("AdminGetAllTestRunCounts called")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	results := dataio.ReadResultsStore()

	// If there are no results, return
	if len(results) == 0 {
		log.Print("No results exist; cannot collect stats")
		w.Write([]byte("{}"))
		return
	}

	// If the first result doesn't have a Name, something is very wrong
	if results[0].Name == "" {
		log.Fatalln("First result object was nil, or its name was nil :", results[0])
		w.Write([]byte("{}"))
		return
	}

	var testruns []string // set of test run IDs
	var stats []structs.Stat

	// Initialise testrunstats with the first result
	stats = append(stats, structs.Stat{TestRunName: results[0].TestRunIdentifier, Count: 0})

	for _, r := range results {
		incremented := false
		// log.Print(":___:", r.TestRunIdentifier, ":___:", stats[0].TestRunName, ":___:")
		if !help.Contains(testruns, r.TestRunIdentifier) {
			testruns = append(testruns, r.TestRunIdentifier)
		}

		// Now collect
		for i := range stats {
			if stats[i].TestRunName == r.TestRunIdentifier {
				stats[i].Count++ // this woun't increment if we loop by object, only by index
				incremented = true
				break
			}
		}
		if !incremented {
			stats = append(stats, structs.Stat{TestRunName: r.TestRunIdentifier, Count: 1})
		}

	}

	// log.Print("Stats of all test runs : ", stats)

	message, _ := json.Marshal(stats)
	// message := stats.ToJSON() // stats is an array, but ToJSON() is on the object
	messageBytes := bytes.TrimSpace([]byte(message))
	w.Write(messageBytes)
	// w.Write([]byte(strconv.Itoa(count)))
}
