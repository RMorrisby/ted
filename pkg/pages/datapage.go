package pages

import (
	"bytes"
	"encoding/json"
	"net/http"
	"ted/pkg/constants"
	"ted/pkg/dataio"
	_ "ted/pkg/handler" // TODO enable
	"ted/pkg/help"
	"ted/pkg/structs"
	"time"

	log "github.com/romana/rlog"
)

func DataPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()                      // find the time right now
	DataPageVars := structs.PageVariables{ //store the date and time in a struct
		Date:        now.Format(constants.LayoutDateISO),
		Time:        now.Format(constants.LayoutTimeISO),
		HostAndPort: help.GetHostAndPortExplicit(),
	}

	// ws.ServeWs(ws.WSHub, w, r)

	err := Templates.ExecuteTemplate(w, "data.html", DataPageVars) //execute the template and pass it the struct to fill in the gaps

	if err != nil {
		log.Debug("template executing error: ", err)
	}
}

func DataGetAllResults(w http.ResponseWriter, r *http.Request) {

	log.Debug("DataGetAllResults called")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	results := dataio.ReadResultStore()

	log.Debug("Total result count : ", len(results))

	// TODO sort the results
	// - group them by version ID
	// - TODO check that the version IDs are in the right order by comparing the timestamps
	// - (maybe the version IDs haven't been written in a sortable way)

	message, _ := json.Marshal(results)
	// message := results.ToJSON() // stats is an array, but ToJSON() is on the object
	messageBytes := bytes.TrimSpace([]byte(message))
	w.Write(messageBytes)

}
