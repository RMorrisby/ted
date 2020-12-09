package main

import (
	"encoding/json"
	"fmt"
	"html/template"

	// "io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	// TODO how to import local file??????
	// "github.com/RMorrisby/ted/handler"
	// "ted/handler"
)

const (
	layoutDateISO = "2006-01-02"
	layoutTimeISO = "15:04:05"
)

type PageVariables struct {
	Date string
	Time string
	Port string
}

type result_struct struct {
	Name string
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return p
	}
	return "8080"
}

func getPortWithColon() string {
	return ":" + getPort()
}

func main() {
	http.HandleFunc("/", IndexPage)
	http.HandleFunc("/is-alive", IsAliveHandler)
	http.HandleFunc("/result", ResultHandler)
	log.Fatal(http.ListenAndServe(getPortWithColon(), nil))
}

func IndexPage(w http.ResponseWriter, r *http.Request) {

	now := time.Now()               // find the time right now
	IndexPageVars := PageVariables{ //store the date and time in a struct
		Date: now.Format(layoutDateISO),
		Time: now.Format(layoutTimeISO),
		Port: getPort(),
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

// Handles the /isalive GET request path, returning a simple JSON object
func IsAliveHandler(w http.ResponseWriter, r *http.Request) {

	log.Print("Is-Alive called")

	data := "{\"is-alive\": true}"

	fmt.Fprintf(w, data)
}

// Handles the /result POST request path for receiving new test results
func ResultHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("/result called")
	switch r.Method {
	case "POST":
		// body, err := ioutil.ReadAll(r.Body)
		// if err != nil {
		// 	panic(err)
		// }

		// Now try to parse the POST body from JSON
		var result result_struct
		// err = json.Unmarshal(body, &result)
		// if err != nil {
		// 	panic(err)
		// }
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields() // catch unwanted fields

		err := d.Decode(&result)
		if err != nil {
			// bad JSON or unrecognized json field
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 'name' field is mandatory
		if result.Name == "" {
			http.Error(w, "Missing field 'name' from JSON object", http.StatusBadRequest)
			return
		}

		log.Println("Result received for test", result.Name)
	default:

		log.Println(r.Method, "/result called")

		fmt.Fprintf(w, "Only POST is supported for /result")
	}
}
