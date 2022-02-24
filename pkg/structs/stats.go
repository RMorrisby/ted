package structs

import (
	"encoding/json"
	"fmt"
)

type Stats struct {
	TestRunName   string
	Total         int
	Passed        int
	PassedOnRerun int
	Failed        int
	KnownIssue    int
	NotRun        int
	LastRun       string // the end time of the most recently-finished test
}

func (s Stats) ToJSON() string {
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
