package structs

// Struct containing the summary history of a test run
// Each suite contains (indirectly) results for various test runs
// Each test run comprises one or more tests, with one result per test per test run
// (this matrix may be sparse - not all tests will have existed for every test run)
// Each result will have a test run identifier (e.g. a version number), a status and (optionally) some notes
type HistoryTestRun struct {
	TestRunName     string
	ResultList      []Result
	TotalCount      int // the total number of executed tests in this test run
	SuccessCount    int // the number of test successes in this test run
	FailCount       int // the number of test failures in this test run
	KnownIssueCount int // the number of tests with known issues in this test run
	NotRunCount     int // the number of tests not yet run in this test run
}
