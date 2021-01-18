package help

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
