package structs

import (
	"encoding/json"
	"fmt"
)

// Unlike the Result struct, this struct's purpose is to bring to the results-table UI the data it needs.
// In effect, this is the Result struct with additional values from the Test and Suite it is linked to.

type ResultForUI struct {
	Categories        string
	Dir               string
	TestName          string
	TestRunIdentifier string
	Status            string
	Priority          int
	StartTimestamp    string
	EndTimestamp      string
	RanBy             string
	Message           string
	TedStatus         string
	TedNotes          string
}

func (r ResultForUI) ToJSON() string {
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
