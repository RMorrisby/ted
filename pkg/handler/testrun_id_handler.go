package handler

import (
	"fmt"
	"net/http"
	"ted/pkg/dataio"
	"ted/pkg/help"

	"github.com/huandu/xstrings"
	log "github.com/romana/rlog"
)

// TestRunIDLatestHandler handles the /testrunid/latest GET request path for getting the latest test run ID
func TestRunIDLatestHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("/testrunid/latest called")
	switch r.Method {
	case "GET":
		log.Println(r.Method, "GET /testrunid/latest called")

		data := dataio.GetLatestTestRun()
		fmt.Fprintf(w, data)

	default:
		http.Error(w, "Only GET is supported for /testrunid/latest", http.StatusMethodNotAllowed)
	}
}

// TestRunIDNextHandler handles the /testrunid/next GET request path for getting the next test run ID
func TestRunIDNextHandler(w http.ResponseWriter, r *http.Request) {
	help.LogNewAPICall("/testrunid/next")
	switch r.Method {
	case "GET":
		log.Println(r.Method, "GET /testrunid/next called")

		latest := dataio.GetLatestTestRun()
		next := xstrings.Successor(latest)
		log.Debugf("Latest test run : %s; next : %s", latest, next)
		fmt.Fprintf(w, next)

	default:
		http.Error(w, "Only GET is supported for /testrunid/next", http.StatusMethodNotAllowed)
	}
}
