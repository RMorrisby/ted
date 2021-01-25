package structs

// Simple struct representing the update to a test's Known Issue status.
// This also needs the latest test run, so that the respective result can also be updated.
type KnownIssueUpdate struct {
	TestName              string
	TestRun               string
	IsKnownIssue          bool
	KnownIssueDescription string
}

// func (t Test) ToJSON() string {
// 	b, err := json.Marshal(t)
// 	if err != nil {
// 		fmt.Println(err)
// 		return ""
// 	}
// 	return string(b)
// }
