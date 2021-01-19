package pages

import (
	// "bytes"
	// "encoding/json"
	"net/http"
	"ted/pkg/constants"
	// "ted/pkg/dataio"
	"ted/pkg/help"
	"ted/pkg/structs"
	"time"

	log "github.com/romana/rlog"
)

// Page showing the summary historical view of a test suite
func HistoryPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()                      // find the time right now
	HistoryPageVars := structs.PageVariables{ //store the date and time in a struct
		Date:        now.Format(constants.LayoutDateISO),
		Time:        now.Format(constants.LayoutTimeISO),
		HostAndPort: help.GetHostAndPortExplicit(),
	}

	// ws.ServeWs(ws.WSHub, w, r)

	err := Templates.ExecuteTemplate(w, "history.html", HistoryPageVars) //execute the template and pass it the struct to fill in the gaps

	if err != nil {
		log.Debug("template executing error: ", err)
	}
}


// REST endpoint to get the test history for the given suite
// TODO
func GetHistoryForSuite(w http.ResponseWriter, r *http.Request) {

	help.LogNewAPICall("AdminGetAllSuites")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	suites := dataio.ReadAllSuites()
	help.MarshalJSONAndWriteToResponse(suites, w)
}
