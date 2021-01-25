package structs

import (
	"encoding/json"
	"fmt"
)

// Unlike the Result struct, this struct's purpose is to bring to the results-table UI the data it needs.
// In effect, this is the Result struct with additional values from the Test and Suite it is linked to.

type ResultForUI struct {
	Categories        string `json:",omitempty"`
	Dir               string `json:",omitempty"`
	TestName          string `json:",omitempty"`
	TestRunIdentifier string `json:",omitempty"`
	Status            string `json:",omitempty"`
	Priority          int    `json:",omitempty"`
	StartTimestamp    string `json:",omitempty"`
	EndTimestamp      string `json:",omitempty"`
	RanBy             string `json:",omitempty"`
	Message           string `json:",omitempty"`
	TedStatus         string `json:",omitempty"`
	TedNotes          string `json:",omitempty"`
}

func (r ResultForUI) ToJSON() string {
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
