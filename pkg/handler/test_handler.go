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

// TestHandler handles the /test POST request path for receiving new tests
func TestHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.Method, "/test called")
	switch r.Method {
	case "GET":
		log.Println(r.URL)
		log.Println(r.URL.Query())
		log.Println(r.URL.Query().Get("test"))
		log.Println(r.URL.Query().Get("test") != "")

		name := r.URL.Query().Get("test")
		if name == "" {
			// A test name must be supplied
			s := "No test name supplied to " + r.Method + " " + r.URL.RequestURI() + "; URL must be /test?test=___"
			log.Error(s)
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		test := dataio.GetTest(name)
		if test == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Test '" + name + "' is not registered in TED"))
		} else {
			fmt.Fprintf(w, test.ToJSON())
		}

	case "POST":

		// Now try to parse the POST body from JSON
		var test structs.Test
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields() // catch unwanted fields

		err := d.Decode(&test)
		if err != nil {
			// bad JSON or unrecognized json field
			log.Error("Bad JSON or unrecognized json field", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// result = result.Trim()

		// 'name' field is mandatory
		if test.Name == "" {
			http.Error(w, "Missing field 'name' from JSON object", http.StatusBadRequest)
			return
		}

		log.Println("New test received :", test.Name)

		dataio.WriteTestToDBIfNew(test)
		w.WriteHeader(http.StatusCreated) // return a 201

	case "DELETE":
		name := r.URL.Query().Get("test")
		if name == "" {
			// A test name must be supplied
			s := "No test name supplied to " + r.Method + " " + r.URL.RequestURI() + "; URL must be /test?test=___"
			log.Error(s)
			http.Error(w, s, http.StatusBadRequest)
			return
		}

		success, err := dataio.DeleteTest(name)
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else if success {
			w.WriteHeader(http.StatusOK) // return a 200
		} else {
			http.Error(w, "ERROR - not successful but no error returned!", 500)
		}
	default:
		log.Debug(r.Method, "/test called")
		http.Error(w, "Only GET, POST, DELETE are supported for /test", http.StatusMethodNotAllowed)
	}
}