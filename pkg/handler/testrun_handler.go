package handler

import (
	// "encoding/json"
	// "fmt"

	"net/http"
	// "ted/pkg/constants"
	"ted/pkg/dataio"
	// "ted/pkg/help"
	// "ted/pkg/pages"
	// "ted/pkg/structs"
	// "ted/pkg/ws"
	// "time"
	log "github.com/romana/rlog"
)

// TestRunHandler handles the /testrun DELETE request path for deleting test runs
func TestRunHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("/testrun called")
	switch r.Method {
	// case "GET":
	// 	log.Println(r.Method, "GET /suite called")
	// 	log.Println(r.URL)
	// 	log.Println(r.URL.Query())
	// 	log.Println(r.URL.Query().Get("suite"))
	// 	log.Println(r.URL.Query().Get("suite") != "")

	// 	name := r.URL.Query().Get("suite")
	// 	if name == "" {
	// 		// A suite name must be supplied
	// 		s := "No suite name supplied to " + r.Method + " " + r.URL.RequestURI() + "; URL must be /suite?suite=___"
	// 		log.Error(s)
	// 		http.Error(w, s, http.StatusBadRequest)
	// 		return
	// 	}

	// 	suite := dataio.GetSuite(name)
	// 	if suite == nil {
	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write([]byte("Suite '" + name + "' is not registered in TED"))
	// 	} else {
	// 		fmt.Fprintf(w, suite.ToJSON())
	// 	}

	// case "POST":

	// 	// Now try to parse the POST body from JSON
	// 	var suite structs.Suite
	// 	d := json.NewDecoder(r.Body)
	// 	d.DisallowUnknownFields() // catch unwanted fields

	// 	err := d.Decode(&suite)
	// 	if err != nil {
	// 		// bad JSON or unrecognized json field
	// 		log.Error("Bad JSON or unrecognized json field", err)
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}

	// 	// result = result.Trim()

	// 	// 'name' field is mandatory
	// 	if suite.Name == "" {
	// 		http.Error(w, "Missing field 'name' from JSON object", http.StatusBadRequest)
	// 		return
	// 	}

	// 	log.Println("New suite received :", suite.Name)

	// 	dataio.WriteSuiteToDBIfNew(suite)
	// 	w.WriteHeader(http.StatusCreated) // return a 201
	case "DELETE":

		name := r.URL.Query().Get("testrun")
		if name == "" {
			// A test run name must be supplied
			s := "No test run name supplied to " + r.Method + " " + r.URL.RequestURI() + "; URL must be /testrun?testrun=___"
			log.Error(s)
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		success, err := dataio.DeleteTestRun(name)
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else if success {
			w.WriteHeader(http.StatusOK) // return a 200
		} else {
			http.Error(w, "ERROR - not successful but no error returned!", 500)
		}

	default:
		log.Debug(r.Method, "/testrun called")
		// http.Error(w, "Only GET, POST, DELETE are supported for /testrun", http.StatusMethodNotAllowed)
		http.Error(w, "Only DELETE is supported for /testrun", http.StatusMethodNotAllowed)
	}
}
