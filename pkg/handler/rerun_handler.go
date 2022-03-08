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

// RerunHandler handles the /rerun GET request path for requesting the list of failed tests
func RerunHandler(w http.ResponseWriter, r *http.Request) {
	
	help.LogNewAPICall("/reruns")

	switch r.Method {
	case "GET":

		testrun := r.URL.Query().Get("testrun")
		if testrun == "" {
			// A test name must be supplied
			s := "No test run supplied to " + r.Method + " " + r.URL.RequestURI() + "; URL must be /reruns?testrun=___"
			log.Error(s)
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		results := dataio.GetFailedTestsForTestrun(testrun)
		help.MarshalJSONAndWriteToResponse(results, w)

	default:
		log.Println(r.Method, "/reruns called")
		http.Error(w, "Only GET is supported for /reruns", http.StatusMethodNotAllowed)
	}
}
