package pages

import (
	"net/http"
	"ted/pkg/dataio"
	"ted/pkg/help"
	"ted/pkg/structs"

	log "github.com/romana/rlog"
)

func DBTestsPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := Templates.ExecuteTemplate(w, "db_tests.html", structs.PageVariables{})

	if err != nil {
		log.Debug("template executing error: ", err)
	}
}

func DBGetEntireTestTable(w http.ResponseWriter, r *http.Request) {

	log.Debug("DBGetEntireTestTable called")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tests := dataio.ReadAllTests()

	log.Debug("Total test count : ", len(tests))

	// TODO sort the results ???
	// - group them by version ID
	// - TODO check that the version IDs are in the right order by comparing the timestamps
	// - (maybe the version IDs haven't been written in a sortable way)

	help.MarshalJSONAndWriteToResponse(tests, w)
}
