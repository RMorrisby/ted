package pages

import (
	"net/http"
	"ted/pkg/constants"
	"ted/pkg/dataio"
	"ted/pkg/help"
	"ted/pkg/structs"
	"time"

	log "github.com/romana/rlog"
)

func DataPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	testrun := r.URL.Query().Get("testrun")
	if testrun == "" {
		// If no testrun has been specified, default to the latest run
		testrun = dataio.LatestTestRun
	}

	now := time.Now()                      // find the time right now
	DataPageVars := structs.PageVariables{ //store the date and time in a struct
		Date:          now.Format(constants.LayoutDateISO),
		Time:          now.Format(constants.LayoutTimeISO),
		HostAndPort:   help.GetHostAndPortExplicit(),
		LatestTestRun: testrun,
	}

	// ws.ServeWs(ws.WSHub, w, r)

	err := Templates.ExecuteTemplate(w, "data.html", DataPageVars) //execute the template and pass it the struct to fill in the gaps

	if err != nil {
		log.Debug("template executing error: ", err)
	}
}

// Gets all results
// If testrun is supplied as a query parameter then only results for that testrun will be returned
func DataGetAllResults(w http.ResponseWriter, r *http.Request) {

	log.Debug("DataGetAllResults called")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	testrun := r.URL.Query().Get("testrun")
	if testrun == "" {
		// If no testrun has been specified, default to the latest run
		testrun = dataio.LatestTestRun
	}

	results := dataio.ReadAllResultsForUI(testrun)

	log.Debug("Total result count : ", len(results))

	// TODO sort the results ???
	// - group them by version ID
	// - TODO check that the version IDs are in the right order by comparing the timestamps
	// - (maybe the version IDs haven't been written in a sortable way)

	help.MarshalJSONAndWriteToResponse(results, w)
}
