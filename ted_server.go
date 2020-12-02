package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	layoutDateISO = "2006-01-02"
	layoutTimeISO = "15:04:05"
)

type PageVariables struct {
	Date string
	Time string
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}

func main() {
	http.HandleFunc("/", IndexPage)
	log.Fatal(http.ListenAndServe(getPort(), nil))
}

func IndexPage(w http.ResponseWriter, r *http.Request) {

	now := time.Now()               // find the time right now
	IndexPageVars := PageVariables{ //store the date and time in a struct
		Date: now.Format(layoutDateISO),
		Time: now.Format(layoutTimeISO),
		Port: getPort()
	}

	t, err := template.ParseFiles("index.html") // parse the html file index.html

	// if there is an error, log it
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, IndexPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps

	if err != nil {
		log.Print("template executing error: ", err)
	}
}
