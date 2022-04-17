package handler

import (
	"fmt"
	"net/http"
	"ted/pkg/dataio"
	"ted/pkg/help"

	"strings"

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
		// next := xstrings.Successor(latest) // this turns 0.0.9 into 0.1.0, but we want 0.0.10

		// Some very manual unpacking & repacking of the version-string
		// Turn 0.0.9 into 00.00.09, increment, then remove leading zeros
		var a [3]string
		for i, x := range strings.Split(latest, ".") {
			if len(x) == 1 {
				a[i] = "0" + x
			} else {
				a[i] = x
			}
		}
		next := strings.Join(a[:], ".")

		succ := xstrings.Successor(next)

		for i, x := range strings.Split(succ, ".") {
			if string(x[0]) == "0" {
				a[i] = string(x[1])
			} else {
				a[i] = x
			}

		}
		next = strings.Join(a[:], ".") // this is now the next testrun ID

		log.Debugf("Latest test run : %s; next : %s", latest, next)
		fmt.Fprintf(w, next)

	default:
		http.Error(w, "Only GET is supported for /testrunid/next", http.StatusMethodNotAllowed)
	}
}
