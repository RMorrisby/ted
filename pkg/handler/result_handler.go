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

// ResultHandler handles the /result POST request path for receiving new test results
// Also handles the /result PUT request path for receiving test result updates
func ResultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println()
	log.Debug("/result called with", r.Method)

	switch r.Method {
	// POST is for new results, PUT is for reruns/updates
	case "POST", "PUT":
		// Now try to parse the body from JSON
		body := r.Body
		// data, _ := ioutil.ReadAll(body)
		// log.Debug(string(data))
		// log.Debug("Result body received :", body)
		var result structs.Result
		d := json.NewDecoder(body)
		d.DisallowUnknownFields() // catch unwanted fields

		err := d.Decode(&result)
		if err != nil {
			// bad JSON or unrecognized json field
			log.Error("Bad JSON or unrecognized json field", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result = result.Trim()

		// 'name' field is mandatory
		if result.TestName == "" {
			http.Error(w, "Missing field 'TestName' from JSON object", http.StatusBadRequest)
			return
		}

		log.Debug("Result received for test", result.TestName)
		log.Debug(result)

		// If the test is not registered, return an error
		if !dataio.TestExists(result.TestName) {
			s := "Result referred to a test that was not registered"
			log.Error(s)
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		// If the test's suite is not registered, return an error
		if !dataio.SuiteExists(result.SuiteName) {
			s := "Result referred to a suite that was not registered"
			log.Error(s)
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		if result.Overwrite {
			// POST requires no Overwrite flag
			if r.Method == "POST" {
				e := fmt.Sprintf("Result received on POST, but the Overwrite flag was true for test %s for testrun %s", result.TestName, result.TestRunIdentifier)
				log.Error(e)
				http.Error(w, e, http.StatusBadRequest)
				return
			}
		} else {
			// PUT requires the Overwrite flag
			if r.Method == "PUT" {
				e := fmt.Sprintf("Result received on PUT, but the Overwrite flag was false for test %s for testrun %s", result.TestName, result.TestRunIdentifier)
				log.Error(e)
				http.Error(w, e, http.StatusBadRequest)
				return
			}
		}

		// If this is a rerun/update but there is no existing result for this testrun, reject it
		// If the test already has a result for this testrun, and this result is not a rerun/update, reject it
		existingResult := dataio.ReadResult(result.TestName, result.TestRunIdentifier)
		log.Debugf("Reading result :: existingResult != nil? %t :: %s", (existingResult != nil), r.Method)
		if existingResult != nil {
			if r.Method == "POST" {
				e := fmt.Sprintf("Result received on POST, but there was an existing result in the DB for test %s for testrun %s", result.TestName, result.TestRunIdentifier)
				log.Error(e)
				http.Error(w, e, http.StatusBadRequest)
				return
			}
		} else {
			if r.Method == "PUT" {
				e := fmt.Sprintf("Result received on PUT, but there was no existing result in the DB for test %s for testrun %s", result.TestName, result.TestRunIdentifier)
				log.Error(e)
				http.Error(w, e, http.StatusBadRequest)
				return
			}
		}

		log.Debug("Result received :", result.ToJSON())
		// TODO users should NOT supply TedStatus or TedNotes fields
		// Amend Result struct
		if result.TedStatus != "" {
			result.TedStatus = ""
		}
		if result.TedNotes != "" {
			result.TedNotes = ""
		}

		DetermineTEDStatusAndNotesForNewResult(&result)
		log.Debug("Result after amendment :", result.ToJSON())

		// The result has passed validation, so now we can write it to the DB and then return the response

		switch r.Method {
		// POST is for new results, PUT is for reruns/updates
		case "POST":
			dataio.WriteResultToStore(result)
			w.WriteHeader(http.StatusCreated) // return a 201
		case "PUT":
			dataio.WriteResultUpdate(result, existingResult)
			w.WriteHeader(http.StatusOK) // return a 200

			// debugging only
			existingResult2 := dataio.ReadResult(result.TestName, result.TestRunIdentifier)
			log.Debug(":: Result after rerun ::")
			log.Debug(existingResult2)
		}

		// If this result does not belong to the latest test run, update the cached variable
		if result.TestRunIdentifier != dataio.LatestTestRun {
			dataio.LatestTestRun = result.TestRunIdentifier
		}

		// If this result does not belong to the latest suite, update the cached variable
		if result.SuiteName != dataio.LatestSuite {
			dataio.LatestSuite = result.SuiteName
		}

	case "PATCH":
		// Now try to parse the body from JSON
		body := r.Body
		// data, _ := ioutil.ReadAll(body)
		// log.Debug(string(data))
		// log.Debug("Result status update body received :", body)
		var resultStatusPatch structs.ResultStatusPatch
		d := json.NewDecoder(body)
		d.DisallowUnknownFields() // catch unwanted fields

		err := d.Decode(&resultStatusPatch)
		if err != nil {
			// bad JSON or unrecognized json field
			log.Error("Bad JSON or unrecognized json field", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resultStatusPatch = resultStatusPatch.Trim()

		// 'name' field is mandatory
		if resultStatusPatch.TestName == "" {
			http.Error(w, "Missing field 'TestName' from JSON object", http.StatusBadRequest)
			return
		}

		log.Debug("PATCH Result received for test", resultStatusPatch.TestName)
		log.Debug(resultStatusPatch)

		// If the test is not registered, return an error
		if !dataio.TestExists(resultStatusPatch.TestName) {
			s := "Result referred to a test that was not registered"
			log.Error(s)
			http.Error(w, s, http.StatusBadRequest)
			return
		}
		
		// If there is no existing result for this testrun, reject this PATCH request
		existingResult := dataio.ReadResult(resultStatusPatch.TestName, resultStatusPatch.TestRunIdentifier)
		log.Debugf("Reading result :: existingResult != nil? %t :: %s", (existingResult != nil), r.Method)
		if existingResult == nil {
			e := fmt.Sprintf("Result status update received on PATCH, but there was no existing result in the DB for test %s for testrun %s", resultStatusPatch.TestName, resultStatusPatch.TestRunIdentifier)
			log.Error(e)
			http.Error(w, e, http.StatusBadRequest)
			return
		}

		log.Debug("PATCH Result received :", resultStatusPatch.ToJSON())

		// The result has passed validation, so now we can write it to the DB and then return the response

		dataio.WriteResultUpdateOnlyStatus(resultStatusPatch.Status, existingResult.Status, resultStatusPatch.TestName, resultStatusPatch.TestRunIdentifier)
		w.WriteHeader(http.StatusOK) // return a 200

	default:
		log.Println(r.Method, "/result called")
		http.Error(w, "Only POST, PUT is supported for /result", http.StatusMethodNotAllowed)
	}
}

// Sets the result's TED status & TED Notes according to what is stored against the test
func DetermineTEDStatusAndNotesForNewResult(result *structs.Result) {
	test := dataio.GetTest(result.TestName)

	if test.IsKnownIssue {
		result.TedStatus = enums.KnownIssue
		result.TedNotes = test.KnownIssueDescription
	} else {
		result.TedStatus = result.Status
		result.TedNotes = ""
	}

	// TODO detect intermittency
}
