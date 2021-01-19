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

// ResultHandler handles the /result POST request path for receiving new test results
func ResultHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("/result called")
	switch r.Method {
	case "POST":

		// Now try to parse the POST body from JSON
		var result structs.Result
		d := json.NewDecoder(r.Body)
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

		dataio.WriteResultToStore(result)
		w.WriteHeader(http.StatusCreated) // return a 201
	default:
		log.Println(r.Method, "/result called")
		http.Error(w, "Only POST is supported for /result", http.StatusMethodNotAllowed)
	}
}
