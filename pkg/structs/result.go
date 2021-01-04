package structs

import (
	"encoding/json"
	"fmt"
)

type Result struct {
	Name              string
	TestRunIdentifier string
	Category          string
	Status            string
	Timestamp         string
	Message           string
}

func ResultHeader() []string {
	header := []string{
		"Name",
		"TestRunIdentifier",
		"Category",
		"Status",
		"Timestamp",
		"Message",
	}
	return header
}

func NewResult(csvLine []string) *Result {
	r := new(Result)
	r.Name = csvLine[0]
	r.TestRunIdentifier = csvLine[1]
	r.Category = csvLine[2]
	r.Status = csvLine[3]
	r.Timestamp = csvLine[4]
	r.Message = csvLine[5]
	return r
}

func (r Result) ToA() []string {
	resultArray := []string{
		r.Name,
		r.TestRunIdentifier,
		r.Category,
		r.Status,
		r.Timestamp,
		r.Message,
	}
	return resultArray
}

func (r Result) ToJSON() string {
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}