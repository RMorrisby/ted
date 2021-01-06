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

func GetHostAndPort() string {
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
		log.Fatal(message, err)
	}
}
