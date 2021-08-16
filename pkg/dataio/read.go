package dataio

import (
	"fmt"
	"strings"
	"ted/pkg/constants"
	"ted/pkg/enums"
	"ted/pkg/structs"

	log "github.com/romana/rlog"
)

var LatestTestRun = ""
var LatestSuite = ""

// TODO remove this and just call  ReadAllResults()
func ReadResultStore() (results []structs.Result) {
	// if help.IsLocal {
	// 	results = ReadResultCSV()
	// } else {
	results = ReadAllResults()
	// }
	return
}

// func ReadResultCSV() []structs.Result {
// 	log.Debug("Will now read results from file :", constants.ResultCSVFilename)
// 	f, err := os.Open(constants.ResultCSVFilename)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer f.Close()

// 	lines, err := csv.NewReader(f).ReadAll()
// 	if err != nil {
// 		panic(err)
// 	}

// 	help.CheckError("Cannot read from file", err)
// 	size := len(lines)
// 	// log.Printf("Read %d results from file", size)
// 	if size == 0 {
// 		log.Fatal("Results CSV is empty - this should always have the header row")
// 		return nil
// 	}

// 	records := make([]structs.Result, size-1)

// 	// Convert each of the lines to a Result (ignoring the header line)
// 	for i, line := range lines[1:] {
// 		result := structs.NewResult(line)
// 		records[i] = *result // we need the * here
// 	}

// 	// debugging
// 	/*
// 		for _, r := range records {
// 			log.Debug(r.Status)
// 		}*/

// 	return records
// }

// Read all results. This will be sent to the UI, so we need to retrieve the extra information like the test name, etc.,
// which is stored in adjacent tables
// If testrun is supplied (not "") then only results for that testrun will be returned
func ReadAllResultsForUI(testrun string) []structs.ResultForUI {
	log.Debug("Reading results from DB for the UI")
	log.Debug("testrun :", testrun)

	sql := constants.ResultTableSelectAllResultsForUISQL
	// If testrun has been specified, add a WHERE clause to the SQL 
	// and change the ORDER BY so that it sorts by result start-time (so that the results are shown in execution-order)
	if testrun != "" {
		// i := strings.Index(sql, " ORDER BY")
		var i = strings.Index(sql, " ORDER BY")
		before := sql[:i]
		// after := sql[i:]
		// sql = before + " WHERE testrun = '" + testrun + "'" + after
		sql = before + " WHERE testrun = '" + testrun + "' ORDER BY result.start_time ASC" 
	}
	log.Debug("SQL :", sql)
	rows, err := DBConn.Query(sql)
	if err != nil {
		log.Criticalf("Error reading results: %q", err)
	}

	var results []structs.ResultForUI
	for rows.Next() {

		var r structs.ResultForUI
		// Categories 	Dir 	Name 	Test Run 	Status 	Priority 	Start 	End 	Ran By 	Message 	TED Status 	TED Notes

		err = rows.Scan(&r.Categories, &r.Dir, &r.TestName, &r.TestRunIdentifier, &r.Status, &r.Priority, &r.StartTimestamp, &r.EndTimestamp, &r.RanBy, &r.Message, &r.TedStatus, &r.TedNotes)
		if err != nil {
			log.Criticalf("Error reading row into struct: %q", err)
		}

		results = append(results, r)
	}

	log.Debugf("Found %d results in DB", len(results))
	return results
}

// Read a single result from the DB
func ReadResult(testname string, testrun string) *structs.Result {
	log.Debug("Reading result from DB")

	sql := fmt.Sprintf("SELECT suite.name, test.name, result.testrun, result.status, result.start_time, result.end_time, result.ran_by, result.message, result.ted_status, result.ted_notes FROM %s result LEFT JOIN %s suite ON result.suite_id = suite.id LEFT JOIN %s test ON result.test_id = test.id WHERE result.testrun = '%s' AND test.name = '%s'", constants.ResultTable, constants.SuiteTable, constants.RegisteredTestTable, testrun, testname)
	log.Debug("SQL :", sql)

	r := structs.Result{}
	// QueryRow is supposed to return an error if there was no row
	// If there was no error, then there was a row
	err := DBConn.QueryRow(sql).Scan(&r.SuiteName, &r.TestName, &r.TestRunIdentifier, &r.Status, &r.StartTimestamp, &r.EndTimestamp, &r.RanBy, &r.Message, &r.TedStatus, &r.TedNotes)
	if err != nil {
		log.Debugf("Result %s :: %s was not found in the DB", testname, testrun)
		return nil
	}

	log.Debugf("Found result %s :: %s", testname, testrun)
	return &r
}

func ReadAllResults() []structs.Result {
	log.Debug("Reading results from DB")

	sql := constants.ResultTableSelectAllSQL
	log.Debug("SQL :", sql)
	rows, err := DBConn.Query(sql)
	if err != nil {
		log.Criticalf("Error reading results: %q", err)
	}

	var results []structs.Result
	for rows.Next() {

		var r structs.Result
		// var rowID int
		err = rows.Scan(&r.SuiteName, &r.TestName, &r.TestRunIdentifier, &r.Status, &r.StartTimestamp, &r.EndTimestamp, &r.RanBy, &r.Message, &r.TedStatus, &r.TedNotes)
		if err != nil {
			log.Criticalf("Error reading row into struct: %q", err)
		}

		results = append(results, r)
	}

	log.Debugf("Found %d results in DB", len(results))
	return results
}

func ReadAllResultsForSuite(suiteName string) []structs.Result {
	log.Debug("Reading results from DB for suite ", suiteName)

	// "SELECT suite.name, test.name, result.testrun, result.status, result.start_time, result.end_time, result.ran_by,
	//         result.message, result.ted_status, result.ted_notes
	// FROM " + ResultTable + " result
	// LEFT JOIN " + SuiteTable + " suite ON result.suite_id = suite.id
	// LEFT JOIN " + RegisteredTestTable + " test ON result.test_id = test.id"

	sql := fmt.Sprintf("%s WHERE suite.name = '%s' ORDER BY result.testrun ASC, test.name ASC", constants.ResultTableSelectAllNoSortingSQL, suiteName)
	log.Debug("SQL :", sql)
	rows, err := DBConn.Query(sql)
	if err != nil {
		log.Criticalf("Error reading results: %q", err)
	}

	var results []structs.Result
	for rows.Next() {

		var r structs.Result
		// var rowID int
		err = rows.Scan(&r.SuiteName, &r.TestName, &r.TestRunIdentifier, &r.Status, &r.StartTimestamp, &r.EndTimestamp, &r.RanBy, &r.Message, &r.TedStatus, &r.TedNotes)
		if err != nil {
			log.Criticalf("Error reading row into struct: %q", err)
		}

		results = append(results, r)
	}

	log.Debugf("Found %d results in DB for suite %s", len(results), suiteName)
	return results
}

// Get the set of all test runs for the given suite.
// Because tests runs are not linked directly to a suite, we have to get all of the results for the suite
// then get the set of test runs for those results.
func GetAllTestRunsForSuite(suiteName string) []string {
	log.Debug("Reading distinct test runs from DB for suite ", suiteName)

	sql := fmt.Sprintf("%s WHERE suite.name = '%s' ORDER BY result.testrun ASC", constants.ResultTableSelectDistinctTestRunNoSortingSQL, suiteName)
	log.Debug("SQL :", sql)
	rows, err := DBConn.Query(sql)
	if err != nil {
		log.Criticalf("Error reading results: %q", err)
	}

	var testruns []string
	for rows.Next() {

		var s string
		// var rowID int
		err = rows.Scan(&s)
		if err != nil {
			log.Criticalf("Error reading test run string from SQL results: %q", err)
		}

		testruns = append(testruns, s)
	}

	log.Debugf("Found %d test runs in DB for suite %s", len(testruns), suiteName)
	return testruns
}

func ReadAllTests() (tests []structs.Test) {

	log.Debug("Reading tests from DB")

	sql := constants.RegisteredTestTableSelectAllSQL + " ORDER BY dir ASC, name ASC"
	log.Debug("SQL :", sql)
	rows, err := DBConn.Query(sql)
	if err != nil {
		log.Criticalf("Error reading tests: %q", err)
	}
	// "SELECT name, dir, priority, categories, description, notes, owner, is_known_issue, known_issue_description from " + RegisteredTestTable
	for rows.Next() {
		var t structs.Test
		err = rows.Scan(&t.Name, &t.Dir, &t.Priority, &t.Categories, &t.Description, &t.Notes, &t.Owner, &t.IsKnownIssue, &t.KnownIssueDescription)
		if err != nil {
			log.Criticalf("Error reading row into struct: %q", err)
		}

		tests = append(tests, t)
	}

	log.Debugf("Found %d tests in DB", len(tests))
	return tests
}

// func GetFailedTestsForTestrun(testrun string) (tests []structs.Test) {
// 	log.Debug("Reading failed tests from DB for testrun", testrun)

// 	sql := "SELECT suite.name, test.name, result.testrun, result.status, result.start_time, result.end_time, result.ran_by, result.message, result.ted_status, result.ted_notes FROM " + constants.ResultTable + " result LEFT JOIN " + constants.SuiteTable + " suite ON result.suite_id = suite.id LEFT JOIN " + constants.RegisteredTestTable + " test ON result.test_id = test.id WHERE result.testrun = " + testrun + " AND test.name = " + testname
// 	log.Debug("SQL :", sql)
// 	rows, err := DBConn.Query(sql)
// 	if err != nil {
// 		log.Criticalf("Error reading tests: %q", err)
// 	}
// 	// "SELECT name, dir, priority, categories, description, notes, owner, is_known_issue, known_issue_description from " + RegisteredTestTable
// 	for rows.Next() {
// 		var t structs.Test
// 		err = rows.Scan(&t.Name, &t.Dir, &t.Priority, &t.Categories, &t.Description, &t.Notes, &t.Owner, &t.IsKnownIssue, &t.KnownIssueDescription)
// 		if err != nil {
// 			log.Criticalf("Error reading row into struct: %q", err)
// 		}

// 		tests = append(tests, t)
// 	}

// 	log.Debugf("Found %d tests in DB", len(tests))
// 	return tests
// }
// Get all tests for the given testrun that did not pass
func GetFailedTestsForTestrun(testrun string) []structs.ResultForUI {
	log.Debug("Reading failed results from DB for testrun", testrun)

	sql := "SELECT test.dir, test.name FROM " + constants.ResultTable + " result LEFT JOIN " + constants.RegisteredTestTable + " test ON result.test_id = test.id WHERE testrun = '" + testrun + "' AND status != '" + string(enums.Passed) + "' ORDER BY test.name ASC"
	log.Debug("SQL :", sql)
	rows, err := DBConn.Query(sql)
	if err != nil {
		log.Criticalf("Error reading results: %q", err)
	}

	var results []structs.ResultForUI
	for rows.Next() {

		var r structs.ResultForUI
		err = rows.Scan(&r.Dir, &r.TestName)
		if err != nil {
			log.Criticalf("Error reading row into struct: %q", err)
		}

		results = append(results, r)
	}

	log.Debugf("Found %d failed results in DB for testrun %s", len(results), testrun)
	return results
}

func ReadAllSuites() (suites []structs.Suite) {

	log.Debug("Reading suites from DB")

	sql := constants.SuiteTableSelectAllSQL
	log.Debug("SQL :", sql)
	rows, err := DBConn.Query(sql)
	if err != nil {
		log.Criticalf("Error reading suites: %q", err)
	}

	for rows.Next() {
		var s structs.Suite
		err = rows.Scan(&s.Name, &s.Description, &s.Owner, &s.Notes)
		if err != nil {
			log.Criticalf("Error reading row into struct: %q", err)
		}

		suites = append(suites, s)
	}

	log.Debugf("Found %d suites in DB", len(suites))
	return suites
}

// May return nil
func GetSuite(name string) *structs.Suite {
	log.Debug("\n\n")
	log.Printf("Reading suites from DB; want suite '%s'", name)

	sql := constants.SuiteTableSelectAllSQL + " WHERE name = '" + name + "'"
	log.Debug("SQL :", sql)

	suite := structs.Suite{}
	// QueryRow is supposed to return an error if there was no row
	// If there was no error, then there was a row
	err := DBConn.QueryRow(sql).Scan(&suite.Name, &suite.Description, &suite.Owner, &suite.Notes)
	if err != nil {
		log.Printf("Suite %s was not found in the DB", name)
		return nil
	}

	if name != suite.Name {
		log.Criticalf("Suite %s was returned from the DB, when we searched for suite %s", suite.Name, name)
		return nil
	}
	log.Debug("Found suite", name)
	return &suite
}

func SuiteExists(name string) bool {
	suite := GetSuite(name)
	if suite == nil {
		return false
	}
	return true
}

// May return nil
func GetTest(name string) *structs.Test {
	log.Printf("Reading tests from DB; want test '%s'", name)
	// = "SELECT name, dir, priority, categories, description, notes, is_known_issue, known_issue_description from "
	sql := constants.RegisteredTestTableSelectAllSQL + " WHERE test.name = '" + name + "'"
	log.Debug("SQL :", sql)

	test := structs.Test{}
	// QueryRow is supposed to return an error if there was no row
	// If there was no error, then there was a row
	err := DBConn.QueryRow(sql).Scan(&test.Name, &test.Dir, &test.Priority, &test.Categories, &test.Description, &test.Notes, &test.Owner, &test.IsKnownIssue, &test.KnownIssueDescription)
	if err != nil {
		log.Printf("Test %s was not found in the DB", name)
		return nil
	}

	if name != test.Name {
		log.Criticalf("Test %s was returned from the DB, when we searched for test %s", test.Name, name)
		return nil
	}
	log.Debug("Found test", name)
	return &test
}

func TestExists(name string) bool {
	test := GetTest(name)
	if test == nil {
		return false
	}
	return true
}

// For each test name supplied, return a partial structs.Test; this contains the test name, dir and categories, and the known-issue info.
func GetTestSummariesFromNames(names []string) []structs.Test {
	log.Printf("Reading tests from DB; want %d tests", len(names))

	nameListSQL := "'" + strings.Join(names, "', '") + "'"
	sql := "SELECT name, dir, categories, is_known_issue, known_issue_description from " + constants.RegisteredTestTable + " WHERE test.name in (" + nameListSQL + ") ORDER BY name ASC"
	log.Debug("SQL :", sql)

	var tests []structs.Test

	log.Debug("SQL :", sql)
	rows, err := DBConn.Query(sql)
	if err != nil {
		log.Criticalf("Error reading tests: %q", err)
	}
	for rows.Next() {
		var t structs.Test
		err = rows.Scan(&t.Name, &t.Dir, &t.Categories, &t.IsKnownIssue, &t.KnownIssueDescription)
		if err != nil {
			log.Criticalf("Error reading row into struct: %q", err)
		}
		log.Debug(t) // TODO remove
		tests = append(tests, t)
	}

	log.Debugf("Retrieved %d tests in DB", len(tests))
	if len(names) != len(tests) {
		log.Errorf("Wanted summaries of %d tests, but only retrieved %d tests", len(names), len(tests))
	}

	return tests
}

// Return the name of the latest test run. This is determined from the most recent result.
// If users wish to do test runs out-of-order, that's their choice.
// Latest == most recent.
// May return "" (only if there are no results, which is an extremely small edge-case)
func GetLatestTestRun() string {
	log.Debug("Reading latest test run from DB")

	// SQL ideas :
	// SELECT timestamp,value,card FROM my_table WHERE id=(select max(id) from my_table)
	// SELECT id,value,card FROM my_table ORDER BY id DESC LIMIT 1;
	sql := "SELECT testrun FROM " + constants.ResultTable + " WHERE id=(SELECT MAX(id) from " + constants.ResultTable + ")"
	log.Debug("SQL :", sql)

	var latestTestRun string
	// QueryRow is supposed to return an error if there was no row
	// If there was no error, then there was a row
	err := DBConn.QueryRow(sql).Scan(&latestTestRun)
	if err != nil {
		log.Error("Failed to determine the most recent test run from the DB")
		return ""
	}

	log.Debug("Latest test run :", latestTestRun)
	return latestTestRun
}

// Return the name of the suite for the latest result.
// Latest == most recent.
// May return "" (only if there are no results, which is sn extremely small edge-case)
func GetSuiteForLatestResult() string {
	log.Debug("Reading suite for latest result from DB")

	sql := "SELECT suite.name FROM " + constants.ResultTable + " result LEFT JOIN " + constants.SuiteTable + " suite ON result.suite_id = suite.id WHERE result.id=(SELECT MAX(id) from " + constants.ResultTable + ")"
	log.Debug("SQL :", sql)

	var latestSuite string
	// QueryRow is supposed to return an error if there was no row
	// If there was no error, then there was a row
	err := DBConn.QueryRow(sql).Scan(&latestSuite)
	if err != nil {
		log.Error("Failed to determine the suite for the most recent result from the DB")
		return ""
	}

	log.Debug("Suite for latest result :", latestSuite)
	return latestSuite
}
