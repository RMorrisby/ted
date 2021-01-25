package structs

import (
	"encoding/json"
	"fmt"

	"strings"
)

/*
Statuses    : PASSED, FAILED
TedStatuses : PASSED, FAILED, NOT RUN, UNKNOWN, KNOWN ISSUE, INTERMITTENT

NOT RUN == TED has received a registration for the test for the upcoming test run.
			Should be used by users to tell TED which tests they intend to execute as part of a test run.
			When the test result comes in, TedStatus should be overwritten.
UNKNOWN == TED did not receive a registration for the test for the upcoming test run.
			This effectively means that TED has no result for the test for the given test run.
			Should only be used internally, to indicate that there is no known result for a given test & test run,
			to assist with summary & historical analysis.
KNOWN ISSUE == A user has marked this test as being affected by a Known Issue. TedNotes field should contain the
				pertinent info.
INTERMITTENT == TODO A user (or TED?) has marked this test as being inconsistent, and failing intermittently.
				TedNotes field should contain the pertinent info.
*/
type Result struct {
	TestName          string
	SuiteName         string
	TestRunIdentifier string
	Status            string
	StartTimestamp    string
	EndTimestamp      string
	RanBy             string
	Message           string
	TedStatus         string
	TedNotes          string
	Overwrite         bool `json:",omitempty"` // used only for reruns/result updates
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

// TODO alter in-place without returning
func (r Result) Trim() Result {
	r.SuiteName = strings.TrimSpace(r.SuiteName)
	r.TestName = strings.TrimSpace(r.TestName)
	r.TestRunIdentifier = strings.TrimSpace(r.TestRunIdentifier)
	r.Status = strings.TrimSpace(r.Status)
	r.StartTimestamp = strings.TrimSpace(r.StartTimestamp)
	r.EndTimestamp = strings.TrimSpace(r.EndTimestamp)
	r.RanBy = strings.TrimSpace(r.RanBy)
	r.Message = strings.TrimSpace(r.Message)

	return r
}

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

func (r Result) ToJSON() string {
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
