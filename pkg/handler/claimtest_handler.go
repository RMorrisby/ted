package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "ted/pkg/constants"
	"ted/pkg/dataio"
	"ted/pkg/enums"

	// "ted/pkg/help"
	// "ted/pkg/pages"
	"ted/pkg/structs"
	// "ted/pkg/ws"
	// "time"
	log "github.com/romana/rlog"
)

// ClaimTestHandler handles the /claimtest POST request path for a test runner to claim a test
func ClaimTestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println()
	log.Debug("/claimtest called")

	switch r.Method {
	case "POST":
		// Now try to parse the body from JSON
		body := r.Body
		// data, _ := ioutil.ReadAll(body)
		// log.Debug(string(data))
		// log.Debug("Claim body received :", body)
		var claim structs.ClaimTest
		d := json.NewDecoder(body)
		d.DisallowUnknownFields() // catch unwanted fields

		err := d.Decode(&claim)
		if err != nil {
			// bad JSON or unrecognized json field
			log.Error("Bad JSON or unrecognized json field", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claim = claim.Trim()

		// 'name' field is mandatory
		if claim.TestName == "" {
			http.Error(w, "Missing field 'TestName' from JSON object", http.StatusBadRequest)
			return
		}
		// TODO testrun is mandatory? IsRerun?

		log.Debug("Claim-test request received for test", claim.TestName)
		log.Debug(claim)

		// If the test is not registered, return an error
		if !dataio.TestExists(claim.TestName) {
			s := "Claim referred to a test that was not registered"
			log.Error(s)
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		// if result.Overwrite {
		// 	// POST requires no Overwrite flag
		// 	if r.Method == "POST" {
		// 		e := fmt.Sprintf("Result received on POST, but the Overwrite flag was true for test %s for testrun %s", result.TestName, result.TestRunIdentifier)
		// 		log.Error(e)
		// 		http.Error(w, e, http.StatusBadRequest)
		// 		return
		// 	}
		// } else {
		// 	// PUT requires the Overwrite flag
		// 	if r.Method == "PUT" {
		// 		e := fmt.Sprintf("Result received on PUT, but the Overwrite flag was false for test %s for testrun %s", result.TestName, result.TestRunIdentifier)
		// 		log.Error(e)
		// 		http.Error(w, e, http.StatusBadRequest)
		// 		return
		// 	}
		// }

		// If the specified test does not exist for that test run, reject the request
		// TODO allow a claim on a test for a testrun without requiring a preexisting NOT RUN result
		existingResult := dataio.ReadResult(claim.TestName, claim.TestRunIdentifier)
		if existingResult == nil {
			e := fmt.Sprintf("Claim received for test %s for testrun %s, but there was no existing result in the DB to claim", claim.TestName, claim.TestRunIdentifier)
			log.Error(e)
			http.Error(w, e, http.StatusBadRequest)
			return
		}

		// If the existing test result has already been run, reject the request
		// TODO claim test for reruns?
		if existingResult.Status != string(enums.NotRun) {
			e := fmt.Sprintf("Claim received for test %s for testrun %s, but the result has already been run", claim.TestName, claim.TestRunIdentifier)
			log.Info(e)

			w.WriteHeader(http.StatusOK) // return a 200 with a body saying 'false'
			w.Write([]byte("false"))
			return
		}

		// If the existing test result has already been claimed, reject the request
		if existingResult.TedStatus == string(enums.Claimed) {
			e := fmt.Sprintf("Claim received for test %s for testrun %s, but the result has already been claimed", claim.TestName, claim.TestRunIdentifier)
			log.Info(e)

			w.WriteHeader(http.StatusOK) // return a 200 with a body saying 'false'
			w.Write([]byte("false"))
			return
		}

		log.Debug("Claim received :", claim.ToJSON())

		// The result has passed validation, so now we can write it to the DB and then return the response

		// Use dataio.WriteResultUpdate()
		// Therefore we need to construct an "updated" result

		// update := existingResult // dataio.WriteResultUpdate() doesnot like this way of creating the object
		update := existingResult.Trim() // use this as a way of copying the existingResult

		update.TedStatus = string(enums.Claimed) // mark the result as claimed

		dataio.WriteResultUpdate(update, existingResult)
		w.WriteHeader(http.StatusOK) // return a 200 with a body saying 'true'
		w.Write([]byte("true"))

	default:
		log.Println(r.Method, "/claimtest called")
		http.Error(w, "Only POST is supported for /claimtest", http.StatusMethodNotAllowed)
	}
}
