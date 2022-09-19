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


func (t Test) ToJSON() string {
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
