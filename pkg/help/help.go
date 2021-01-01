package help

import (
	"log"
	"os"
)

var IsLocal bool // cache the fact that we are running locally (or not) // should be available globally

// If "PORT" is set, we are not running locally
func IsTEDRunningLocally() bool {
	p := os.Getenv("PORT")
	if p != "" {
		return false
	}
	return true
}

func CheckError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}