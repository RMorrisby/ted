package handler

import (
	"encoding/json"

	"net/http"
	// "ted/pkg/constants"
	"ted/pkg/dataio"
	"ted/pkg/enums"
	"ted/pkg/help"

	// "ted/pkg/pages"
	"ted/pkg/structs"
	// "ted/pkg/ws"
	// "time"
	log "github.com/romana/rlog"
)

// PauseHandler handles the /pause GET, PUT request paths for pauses
func PauseHandler(w http.ResponseWriter, r *http.Request) {
	help.LogNewAPICall("/pause")
	switch r.Method {
	case "GET":
		// TODO require a testrun-id?
		// log.Println(r.URL)
		// log.Println(r.URL.Query())
		// log.Println(r.URL.Query().Get("pause"))
		// log.Println(r.URL.Query().Get("pause") != "")

		// name := r.URL.Query().Get("pause")
		// if name == "" {
		// 	// A pause name must be supplied
		// 	s := "No pause name supplied to " + r.Method + " " + r.URL.RequestURI() + "; URL must be /pause?pause=___"
		// 	log.Error(s)
		// 	http.Error(w, s, http.StatusBadRequest)
		// 	return
		// }

		pause := dataio.GetStatus(enums.StatusNamePause)
		if pause == nil {
			w.WriteHeader(http.StatusOK) // return a 200 with a body saying 'false'
			w.Write([]byte("false"))
		} else if pause.Value == enums.Unpaused {
			w.WriteHeader(http.StatusOK) // return a 200 with a body saying 'false'
			w.Write([]byte("false"))
		} else if pause.Value == enums.Paused {
			w.WriteHeader(http.StatusOK) // return a 200 with a body saying 'true'
			w.Write([]byte("true"))
		} else {
			w.WriteHeader(http.StatusInternalServerError) // return a 500 with a body saying 'unknown'
			w.Write([]byte("unknown"))
		}

	case "PUT":

		// Now try to parse the PUT body from JSON
		var pause structs.Status
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields() // catch unwanted fields

		err := d.Decode(&pause)
		if err != nil {
			// bad JSON or unrecognized json field
			log.Error("Bad JSON or unrecognized json field", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// result = result.Trim()

		// 'name' field is mandatory
		if pause.Name == "" {
			http.Error(w, "Missing field 'name' from JSON object", http.StatusBadRequest)
			return
		}

		// 'type' field is mandatory
		if pause.Type == "" {
			http.Error(w, "Missing field 'type' from JSON object", http.StatusBadRequest)
			return
		}

		// 'value' field is mandatory
		if pause.Value == "" {
			http.Error(w, "Missing field 'value' from JSON object", http.StatusBadRequest)
			return
		}

		log.Println("New pause-request received :", pause.Name, pause.Type, pause.Value)

		pause = help.SanitiseStatus(pause)

		// If it's not a valid pause-request, reject it

		if pause.Type != enums.Pause {
			http.Error(w, "A pause-request must have a 'type' of '"+enums.Pause+"'", 400)
			return
		}

		if pause.Name != enums.StatusNamePause {
			http.Error(w, "A pause-request must have a 'name' of '"+enums.StatusNamePause+"'", 400)
			return
		}

		if pause.Value != enums.Paused && pause.Value != enums.Unpaused {
			http.Error(w, "A pause-request must have a 'value' of '"+enums.Paused+"' or '"+enums.Unpaused+"'", 400)
			return
		}

		// Now write the status-object to the DB, overwriting if necessary
		success, err := dataio.WriteStatusToDB(pause)

		if err != nil {
			http.Error(w, err.Error(), 400)
		} else if success {
			w.WriteHeader(http.StatusOK) // return a 200
		} else {
			http.Error(w, "ERROR - not successful but no error returned!", 500)
		}

	default:
		log.Debug(r.Method, "/pause called")
		http.Error(w, "Only GET, PUT are supported for /pause", http.StatusMethodNotAllowed)
	}
}
