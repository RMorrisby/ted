package structs

// Struct containing the summary history of all tests results belonging to a suite
// Each suite contains (indirectly) results for various test runs
// Each test run comprises one or more tests, with one result per test per test run
// (this matrix may be sparse - not all tests will have existed for every test run)
// Each result will have a test run identifier (e.g. a version number), a status and (optionally) some notes
type HistorySuite struct {
	SuiteName       string
	TestRuns        []string                  // list of all of the test runs the suite has any results for
	TestRunMap      map[string]HistoryTestRun // map of the test run name and its results (Result struct)
	Tests           []Test                    // list of all of the tests the suite has any results for
	TotalCount      int                       // the total number of executed tests in the most recent test run
	SuccessCount    int                       // the number of test successes in the most recent test run
	FailCount       int                       // the number of test failures in the most recent test run
	KnownIssueCount int                       // the number of tests with known issues in the most recent test run
	NotRunCount     int                       // the number of tests not yet run in the most recent test run
}
