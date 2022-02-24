package handler

import (

	// "fmt"

	"net/http"
	// "ted/pkg/constants"
	"ted/pkg/dataio"
	"ted/pkg/help"

	// "ted/pkg/pages"

	// "ted/pkg/ws"
	// "time"
	log "github.com/romana/rlog"
)

// RerunHandler handles the /rerun GET request path for requesting the stats of a testrun
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	
	switch r.Method {
	case "GET":

		testrun := r.URL.Query().Get("testrun")
		log.Debugf("/stats called for testrun %s", testrun)
		if testrun == "" {
			// A test name must be supplied
			s := "No test run supplied to " + r.Method + " " + r.URL.RequestURI() + "; URL must be /stats?testrun=___"
			log.Error(s)
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		stats := dataio.GetStatsForTestrun(testrun)
		help.MarshalJSONAndWriteToResponse(stats, w)

	default:
		log.Println(r.Method, "/stats called")
		http.Error(w, "Only GET is supported for /stats", http.StatusMethodNotAllowed)
	}
}
