package help

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"ted/pkg/structs"

	log "github.com/romana/rlog"
)

var IsLocal bool // cache the fact that we are running locally (or not) // should be available globally

// Common logging helper so that eacn API can log in a more common way
func LogNewAPICall(methodName string) {
	fmt.Println("") // print a new-line to help make the logs be more readable
	log.Debug(methodName, "called")
}

// If "PORT" is set, we are not running locally
func IsTEDRunningLocally() bool {
	p := os.Getenv("PORT")
	if p != "" {
		return false
	}
	return true
}

func GetHostAndPort() string {
	// If "PORT" is set, we are running on Heroku
	// If not set, we are running locally (Win10)
	p := os.Getenv("PORT")

	// If Heroku, do not specify the hostname. Just return the : and the port
	if p != "" {
		return ":" + p
	}

	// If local (Win10), we should specify localhost as the host
	// This stops Win10 from asking about firewall permissions with each new build
	return "localhost:8080"
}

func GetHostAndPortExplicit() string {
	// If "PORT" is set, we are running on Heroku
	// If not set, we are running locally (Win10)
	p := os.Getenv("PORT")

	// If Heroku, do not specify the hostname. Just return the : and the port
	if p != "" {
		return "arcane-ravine-69473.herokuapp.com:" + p
	}

	// If local (Win10), we should specify localhost as the host
	// This stops Win10 from asking about firewall permissions with each new build
	return "localhost:8080"
}

func CheckError(message string, err error) {
	if err != nil {
		log.Critical(message, err)
	}
}

// Contains asks whether the string list contains the supplied string
func Contains(listOfStrings []string, myString string) bool {
	for _, s := range listOfStrings {
		if s == myString {
			return true
		}
	}
	return false
}

// Common method for sending some data to the REST responsean (as JSON)
func MarshalJSONAndWriteToResponse(obj interface{}, w http.ResponseWriter) {

	message, _ := json.Marshal(obj)
	messageBytes := bytes.TrimSpace([]byte(message))
	w.Write(messageBytes)
}

// Takes a Result and a Test, and forms the matching ResultForUI object
func FormResultForUI(result structs.Result, test *structs.Test) (resultForUI structs.ResultForUI) {

	resultForUI.TestName = result.TestName
	resultForUI.TestRunIdentifier = result.TestRunIdentifier
	resultForUI.Status = result.Status
	resultForUI.StartTimestamp = result.StartTimestamp
	resultForUI.EndTimestamp = result.EndTimestamp
	resultForUI.RanBy = result.RanBy
	resultForUI.Message = result.Message
	resultForUI.TedStatus = result.TedStatus
	resultForUI.TedNotes = result.TedNotes

	resultForUI.Categories = test.Categories
	resultForUI.Dir = test.Dir
	resultForUI.Priority = test.Priority
	return
}

// Takes an array of test run names (e.g. v1.10.2, v1.11.11, v1.9.9) and sorts them properly
// This sorts in-place
func SortTestRuns(testruns []string) {

	sort.SliceStable(testruns, func(i, k int) bool {
		r := regexp.MustCompile(`(\d+).(\d+).(\d+)`)
		rs := r.FindStringSubmatch(testruns[i])

		major1, _ := strconv.Atoi(rs[1])
		minor1, _ := strconv.Atoi(rs[2])
		patch1, _ := strconv.Atoi(rs[3])

		rs = r.FindStringSubmatch(testruns[k])

		major2, _ := strconv.Atoi(rs[1])
		minor2, _ := strconv.Atoi(rs[2])
		patch2, _ := strconv.Atoi(rs[3])

		// Sort by Major
		if major1 < major2 {
			return true
		}
		if major1 > major2 {
			return false
		}

		// Sort by Minor
		if minor1 < minor2 {
			return true
		}
		if minor1 > minor2 {
			return false
		}

		// Sort by Patch
		if patch1 < patch2 {
			return true
		}
		if patch1 > patch2 {
			return false
		}

		return true
	})
}

// replacer replaces ' with \'
// Actually, PostGres likes '' instead of \'
// It's a package-level variable so we can easily reuse it, but
// this program doesn't take advantage of that fact.
// var replacer = strings.NewReplacer("'", "\\'")
var replacer = strings.NewReplacer("'", "''")

// Sanitises the fields within the test object so that SQL-injection can't occur
// TODO better to sanitise each SQL line directly before execution
func SanitiseTest(test structs.Test) structs.Test {
	test.Description = replacer.Replace(test.Description)

	return test
}

// Sanitises the fields within the update object so that SQL-injection can't occur
// TODO better to sanitise each SQL line directly before execution
func SanitiseUpdate(update structs.KnownIssueUpdate) structs.KnownIssueUpdate {
	update.KnownIssueDescription = replacer.Replace(update.KnownIssueDescription)

	return update
}

// Takes the name (perhaps qualified) of a DB column and returns the SQL string for the COALESCE command
// which will format the column value into a date-string or an empty string (if NULL)
// Golang can't automatically convert NULL DB timestamps into strings, so we use COALESCE to get Postgres
// to do the conversion from timestamp (or NULL) into string
func CoalesceDateSQL(column string) string {
	return "COALESCE(to_char(" + column + ", 'YYYY-MM-DD HH24:MI:SS'), '')"
}
