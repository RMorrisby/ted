package structs

import (
	"encoding/json"
	"fmt"
)

type Test struct {
	Name                  string
	Dir                   string
	Priority              int
	Categories            string // pipe-separated string
	Description           string
	Notes                 string
	Owner                 string
	IsKnownIssue          bool
	KnownIssueDescription string
}

// func ResultHeader() []string {
// 	header := []string{
// 		"Name",
// 		"TestRunIdentifier",
// 		"Category",
// 		"Status",
// 		"Timestamp",
// 		"Message",
// 	}
// 	return header
// }

// func NewResult(csvLine []string) *Result {
// 	r := new(Result)
// 	r.Name = csvLine[0]
// 	r.TestRunIdentifier = csvLine[1]
// 	r.Category = csvLine[2]
// 	r.Status = csvLine[3]
// 	r.Timestamp = csvLine[4]
// 	r.Message = csvLine[5]
// 	return r
// }

// // TODO alter in-place without returning
// func (r Result) Trim() Result {
// 	r.SuiteName = strings.TrimSpace(r.SuiteName)
// 	r.Name = strings.TrimSpace(r.Name)
// 	r.TestRunIdentifier = strings.TrimSpace(r.TestRunIdentifier)
// 	r.Status = strings.TrimSpace(r.Status)
// 	r.StartTimestamp = strings.TrimSpace(r.StartTimestamp)
// 	r.EndTimestamp = strings.TrimSpace(r.EndTimestamp)
// 	r.RanBy = strings.TrimSpace(r.RanBy)
// 	r.Message = strings.TrimSpace(r.Message)

// 	return r
// }

// func (r Result) ToA() []string {
// 	resultArray := []string{
// 		r.Name,
// 		r.TestRunIdentifier,
// 		r.Category,
// 		r.Status,
// 		r.Timestamp,
// 		r.Message,
// 	}
// 	return resultArray
// }

func (t Test) ToJSON() string {
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
