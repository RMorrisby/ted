package pages

import (
_	"encoding/json"
_	"fmt"
	"html/template"
_	"path/filepath"

_	"database/sql"

_	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"

	"log"
	"net/http"
_	"os"
	_ "ted/pkg/handler" // TODO enable
	"ted/pkg/structs"
	"ted/pkg/constants"
	"ted/pkg/dataio"
	"time"
)


func DataPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()                       // find the time right now
	DataPageVars := structs.PageVariables{ //store the date and time in a struct
		Date:         now.Format(constants.LayoutDateISO),
		Time:         now.Format(constants.LayoutTimeISO),
	}

	t, err := template.ParseFiles("data.html") // parse the html file index.html

	// if there is an error, log it
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, DataPageVars) //execute the template and pass it the struct to fill in the gaps

	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func DataPage2(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	t, err := template.ParseFiles("data2.html") // parse the html file index.html

	// if there is an error, log it
	if err != nil {
		log.Print("template parsing error: ", err)
	}

 // (results []structs.Result)
results := dataio.ReadResultsStore()

	err = t.Execute(w, results) //execute the template and pass it the struct to fill in the gaps

	if err != nil {
		log.Print("template executing error: ", err)
	}
}
