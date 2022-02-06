package structs

import (
	"encoding/json"
	"fmt"

	"strings"
)

/*
Statuses    : PASSED, FAILED
TedStatuses : PASSED, FAILED, NOT RUN, CLAIMED, UNKNOWN, KNOWN ISSUE, INTERMITTENT

NOT RUN == TED has received a registration for the test for the upcoming test run.
			Should be used by users to tell TED which tests they intend to execute as part of a test run.
			When the test result comes in, TedStatus should be overwritten.
CLAIMED == TED has received a request from a test runner to run this test.
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
type ClaimTest struct {
	TestName          string
	TestRunIdentifier string
	IsRerun           bool `json:",omitempty"` // used only for reruns
}

// TODO alter in-place without returning
func (r ClaimTest) Trim() ClaimTest {
	r.TestName = strings.TrimSpace(r.TestName)
	r.TestRunIdentifier = strings.TrimSpace(r.TestRunIdentifier)

	return r
}

func (r ClaimTest) ToJSON() string {
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
