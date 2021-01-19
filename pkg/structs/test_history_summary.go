package structs

// Struct containing the summary history of a test
// Each suite contains tests (this relationship is not strictly enforced, to allow tests to belong to multiiple suites)
// Each test will have one or more results
// Each result will have a test run identifier (e.g. a version number), a status and (optionally) some notes 
type TestHistorySummary struct {
	TestName         string
	ResultMap map[string]Result
	TotalCount    int
	SuccessCount int
	FailCount    int
	NotRunCount    int
	KnownIssueCount    int
}
