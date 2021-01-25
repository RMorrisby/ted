package handler

import (
	"encoding/json"
	// "fmt"

	"net/http"
	// "ted/pkg/constants"
	"ted/pkg/dataio"
	// "ted/pkg/help"
	// "ted/pkg/pages"
	"ted/pkg/structs"
	// "ted/pkg/ws"
	// "time"
	log "github.com/romana/rlog"
)

// TestUpdateHandler handles the /testupdate POST request path for updating existing tests
func TestUpdateHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("/testupdate called")
	switch r.Method {
	case "POST":

		// Now try to parse the POST body from JSON
		var update structs.KnownIssueUpdate
		log.Debug("/testupdate POST ::", r.Body)
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields() // catch unwanted fields

		err := d.Decode(&update)
		if err != nil {
			// bad JSON or unrecognized json field
			log.Error("Bad JSON or unrecognized json field", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 'name' field is mandatory
		if update.TestName == "" {
			http.Error(w, "Missing field 'Name' from JSON object", http.StatusBadRequest)
			return
		}

		log.Debug("Update received for test", update.TestName)
		log.Debug(update)

		// If the test is not registered, return an error
		if !dataio.TestExists(update.TestName) {
			s := "Update referred to a test that was not registered"
			log.Error(s)
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		// The test has passed validation, so now we can the update to the DB and then return the response
		dataio.WriteTestKnownIssueUpdate(update)

		// Also update the result for the given testrun
		dataio.WriteResultKnownIssueUpdate(update)

		// // If this test does not belong to the latest test run, update the cached variable
		// if test.TestRunIdentifier != dataio.LatestTestRun {
		// 	dataio.LatestTestRun = test.TestRunIdentifier
		// }

		// // If this test does not belong to the latest suite, update the cached variable
		// if test.SuiteName != dataio.LatestSuite {
		// 	dataio.LatestSuite = test.SuiteName
		// }

		// If we are setting or changing the Known Issue, return a 201 (to indicate that something has been written)
		// If we are clearing the Known Issue, return a 200
		if update.IsKnownIssue {
			w.WriteHeader(http.StatusCreated) // return a 201
		} else {
			w.WriteHeader(http.StatusOK) // return a 200
		}
	default:
		log.Println(r.Method, "/testupdate called")
		http.Error(w, "Only POST is supported for /testupdate", http.StatusMethodNotAllowed)
	}
}
