package structs

// Struct containing the summary history of all tests results belonging to a suite
// Each suite contains tests (this relationship is not strictly enforced, to allow tests to belong to multiiple suites)
// Each test will have one or more results
// Each result will have a test run identifier (e.g. a version number), a status and (optionally) some notes 
type History struct {
	SuiteName         string
	TestMap map[string]TestHistorySummary
	TotalCount    int
	SuccessCount int
	FailCount    int
	NotRunCount    int
	KnownIssueCount    int
}
