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
	"strconv"
	"ted/pkg/constants"
	"ted/pkg/dataio"
	_ "ted/pkg/handler" // TODO enable
	"ted/pkg/structs"
	_ "ted/pkg/ws"
	"time"

	_ "github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

// The Admin page
func AdminPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()                       // find the time right now
	AdminPageVars := structs.PageVariables{ //store the date and time in a struct
		Date: now.Format(constants.LayoutDateISO),
		Time: now.Format(constants.LayoutTimeISO),
	}

	// ws.ServeWs(ws.WSHub, w, r)

	err := Templates.ExecuteTemplate(w, "admin.html", AdminPageVars) //execute the template and pass it the struct to fill in the gaps

	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// REST endpoint to trigger the deletion of all results
func AdminDeleteAll(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	success, err := dataio.DeleteAllResults()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	} else if success {
		w.Write([]byte(strconv.Itoa(0)))
	} else {
		http.Error(w, "ERROR - not successful but no error returned!", 500)
		return
	}
}

// REST endpoint to get the total number of results
func AdminGetCount(w http.ResponseWriter, r *http.Request) {

	log.Print("AdminGetCount called")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count := len(dataio.ReadResultsStore())
	log.Print("Total result count : ", count)
	w.Write([]byte(strconv.Itoa(count)))
}
