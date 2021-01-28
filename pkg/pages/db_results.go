package pages

import (
	"net/http"
	"ted/pkg/dataio"
	"ted/pkg/help"
	"ted/pkg/structs"

	log "github.com/romana/rlog"
)

func DBResultsPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := Templates.ExecuteTemplate(w, "db_results.html", structs.PageVariables{})

	if err != nil {
		log.Debug("template executing error: ", err)
	}
}

func DBGetEntireResultTable(w http.ResponseWriter, r *http.Request) {

	log.Debug("DBGetEntireResultTable called")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	results := dataio.ReadAllResults()

	log.Debug("Total result count : ", len(results))

	// TODO sort the results ???
	// - group them by version ID
	// - TODO check that the version IDs are in the right order by comparing the timestamps
	// - (maybe the version IDs haven't been written in a sortable way)

	help.MarshalJSONAndWriteToResponse(results, w)
}
