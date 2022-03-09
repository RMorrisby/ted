package handler

import (
	"encoding/json"
	"fmt"

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
			w.WriteHeader(http.StatusNoContent) // return a 204
			w.Write([]byte("No Pause status is set"))
		} else {
			w.WriteHeader(http.StatusOK) // return a 200
			fmt.Fprintf(w, pause.ToJSON())
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

		success, err := dataio.WritePauseToDBIfNew(pause)

		if err != nil {
			http.Error(w, err.Error(), 400)
		} else if success {
			w.WriteHeader(http.StatusCreated) // return a 201
		} else {
			http.Error(w, "ERROR - not successful but no error returned!", 500)
		}

	// case "DELETE":
	// 	name := r.URL.Query().Get("pause")
	// 	if name == "" {
	// 		// A pause name must be supplied
	// 		s := "No pause name supplied to " + r.Method + " " + r.URL.RequestURI() + "; URL must be /pause?pause=___"
	// 		log.Error(s)
	// 		http.Error(w, s, http.StatusBadRequest)
	// 		return
	// 	}

	// 	success, err := dataio.DeletePause(name)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), 500)
	// 	} else if success {
	// 		w.WriteHeader(http.StatusOK) // return a 200
	// 	} else {
	// 		http.Error(w, "ERROR - not successful but no error returned!", 500)
	// 	}
	default:
		log.Debug(r.Method, "/pause called")
		http.Error(w, "Only GET, PUT are supported for /pause", http.StatusMethodNotAllowed)
	}
}
