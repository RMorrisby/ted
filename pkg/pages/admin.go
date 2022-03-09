package pages

import (
	_ "fmt"
	"net/http"
	"strconv"
	"ted/pkg/constants"
	"ted/pkg/dataio"
	"ted/pkg/help"
	"ted/pkg/structs"
	"time"

	log "github.com/romana/rlog"
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
		log.Debug("template executing error: ", err)
	}
}

// REST endpoint to trigger the deletion of all results
func AdminDeleteAllResults(w http.ResponseWriter, r *http.Request) {

	if r.Method != "DELETE" {
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

// REST endpoint to trigger the deletion of all tests
func AdminDeleteAllTests(w http.ResponseWriter, r *http.Request) {

	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	success, err := dataio.DeleteAllTests()
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

// REST endpoint to trigger the deletion of all suites
func AdminDeleteAllSuites(w http.ResponseWriter, r *http.Request) {

	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	success, err := dataio.DeleteAllSuites()
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

// REST endpoint to trigger the deletion of all statuses
func AdminDeleteAllStatuses(w http.ResponseWriter, r *http.Request) {

	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	success, err := dataio.DeleteAllStatuses()
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
func AdminGetResultCount(w http.ResponseWriter, r *http.Request) {

	help.LogNewAPICall("AdminGetResultCount")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count := len(dataio.ReadAllResults())
	log.Debug("Total result count : ", count)
	w.Write([]byte(strconv.Itoa(count)))
}

// REST endpoint to get all known tests
func AdminGetAllTests(w http.ResponseWriter, r *http.Request) {

	help.LogNewAPICall("AdminGetAllTests")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tests := dataio.ReadAllTests()
	help.MarshalJSONAndWriteToResponse(tests, w)
}

// REST endpoint to get the total number of tests
func AdminGetTestCount(w http.ResponseWriter, r *http.Request) {

	help.LogNewAPICall("AdminGetTestCount")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count := len(dataio.ReadAllTests())
	log.Debug("Total test count : ", count)
	w.Write([]byte(strconv.Itoa(count)))
}

// REST endpoint to get all known suites
func AdminGetAllSuites(w http.ResponseWriter, r *http.Request) {

	help.LogNewAPICall("AdminGetAllSuites")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	suites := dataio.ReadAllSuites()
	help.MarshalJSONAndWriteToResponse(suites, w)
}

// REST endpoint to get all known statuses
func AdminGetAllStatuses(w http.ResponseWriter, r *http.Request) {

	help.LogNewAPICall("AdminGetAllStatuses")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	suites := dataio.ReadAllStatuses()
	help.MarshalJSONAndWriteToResponse(suites, w)
}

// REST endpoint to get the total number of suites
func AdminGetSuiteCount(w http.ResponseWriter, r *http.Request) {

	help.LogNewAPICall("AdminGetSuiteCount")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count := len(dataio.ReadAllSuites())
	log.Debug("Total suite count : ", count)
	w.Write([]byte(strconv.Itoa(count)))
}

// REST endpoint to get the total number of statuses
func AdminGetStatusCount(w http.ResponseWriter, r *http.Request) {

	help.LogNewAPICall("AdminGetStatusCount")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count := len(dataio.ReadAllStatuses())
	log.Debug("Total status count : ", count)
	w.Write([]byte(strconv.Itoa(count)))
}

// REST endpoint to get the total number of known test runs, and the number of results in each run
func AdminGetAllTestRunCounts(w http.ResponseWriter, r *http.Request) {

	help.LogNewAPICall("AdminGetAllTestRunCounts")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	results := dataio.ReadResultStore()

	// If there are no results, return
	if len(results) == 0 {
		log.Debug("No results exist; cannot collect stats")
		w.Write([]byte("{}"))
		return
	}

	// If the first result doesn't have a Name, something is very wrong
	if results[0].TestName == "" {
		log.Critical("First result object was nil, or its name was nil :", results[0])
		w.Write([]byte("{}"))
		return
	}

	var testruns []string // set of test run IDs
	var stats []structs.Stats

	// Initialise testrunstats with the first result
	stats = append(stats, structs.Stats{TestRunName: results[0].TestRunIdentifier, Total: 0})

	for _, r := range results {
		incremented := false
		// log.Debug(":___:", r.TestRunIdentifier, ":___:", stats[0].TestRunName, ":___:")
		if !help.Contains(testruns, r.TestRunIdentifier) {
			testruns = append(testruns, r.TestRunIdentifier)
		}

		// Now collect
		for i := range stats {
			if stats[i].TestRunName == r.TestRunIdentifier {
				stats[i].Total++ // this woun't increment if we loop by object, only by index
				incremented = true
				break
			}
		}
		if !incremented {
			stats = append(stats, structs.Stats{TestRunName: r.TestRunIdentifier, Total: 1})
		}

	}

	// log.Debug("Stats of all test runs : ", stats)
	help.MarshalJSONAndWriteToResponse(stats, w)
}
