package structs

import (
	"encoding/json"
	"fmt"

	"strings"
)

/*
A simple class holding a partial test result update. Only updates the result's status (Passed or Failed).
Should only be received via PATCH.

Statuses    : PASSED, FAILED
*/
type ResultStatusPatch struct {
	TestName          string
	TestRunIdentifier string
	Status            string
}

// TODO alter in-place without returning
func (r ResultStatusPatch) Trim() ResultStatusPatch {
	r.TestName = strings.TrimSpace(r.TestName)
	r.TestRunIdentifier = strings.TrimSpace(r.TestRunIdentifier)
	r.Status = strings.TrimSpace(r.Status)

	return r
}

func (r ResultStatusPatch) ToJSON() string {
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
