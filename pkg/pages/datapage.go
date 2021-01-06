package pages

import (
	_ "database/sql"
	_ "encoding/json"
	_ "fmt"
	_ "html/template"
	"log"
	"net/http"
	_ "os"
	_ "path/filepath"
	"ted/pkg/constants"
	"ted/pkg/dataio"
	_ "ted/pkg/handler" // TODO enable
	"ted/pkg/help"
	"ted/pkg/structs"
	_ "ted/pkg/ws"
	"time"

	_ "github.com/gorilla/websocket"
	_ "github.com/lib/pq"
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
		log.Print("template executing error: ", err)
	}
}

func DataPage2(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// (results []structs.Result)
	results := dataio.ReadResultsStore()

	err := Templates.ExecuteTemplate(w, "data2.html", results) //execute the template and pass it the struct to fill in the gaps

	if err != nil {
		log.Print("template executing error: ", err)
	}
}
