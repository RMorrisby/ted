package dataio

import (
	"fmt"
	"strings"
	"ted/pkg/constants"
	"ted/pkg/structs"

	log "github.com/romana/rlog"
)

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
// 	log.Println("Will now read results from file :", constants.ResultCSVFilename)
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
// 			log.Println(r.Status)
// 		}*/

// 	return records
// }

func ReadAllResultsForUI() []structs.ResultForUI {
	log.Println("Reading results from DB for the UI")

	sql := constants.ResultTableSelectAllResultsForUISQL
	log.Println("SQL :", sql)
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

func ReadAllResults() []structs.Result {
	log.Println("Reading results from DB")

	sql := constants.ResultTableSelectAllSQL
	log.Println("SQL :", sql)
	rows, err := DBConn.Query(sql)
	if err != nil {
		log.Criticalf("Error reading results: %q", err)
	}

	// cols, _ := rows.Columns()
	// log.Printf("Found %d columns in DB", len(cols))
	// log.Printf("Found %d results in DB", resultCount)

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
	log.Println("Reading results from DB for suite ", suiteName)

	// "SELECT suite.name, test.name, result.testrun, result.status, result.start_time, result.end_time, result.ran_by,
	//         result.message, result.ted_status, result.ted_notes
	// FROM " + ResultTable + " result
	// LEFT JOIN " + SuiteTable + " suite ON result.suite_id = suite.id
	// LEFT JOIN " + RegisteredTestTable + " test ON result.test_id = test.id"

	sql := fmt.Sprintf("%s WHERE suite.name = '%s' ORDER BY result.testrun ASC, test.name ASC", constants.ResultTableSelectAllNoSortingSQL, suiteName)
	log.Println("SQL :", sql)
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
	log.Println("Reading distinct test runs from DB for suite ", suiteName)

	sql := fmt.Sprintf("%s WHERE suite.name = '%s' ORDER BY result.testrun ASC", constants.ResultTableSelectDistinctTestRunNoSortingSQL, suiteName)
	log.Println("SQL :", sql)
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

	sql := constants.RegisteredTestTableSelectAllSQL
	log.Println("SQL :", sql)
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

func ReadAllSuites() (suites []structs.Suite) {

	log.Debug("Reading suites from DB")

	sql := constants.SuiteTableSelectAllSQL
	log.Println("SQL :", sql)
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
	log.Println("SQL :", sql)

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
	log.Println("SQL :", sql)

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

// For each test name supplied, return a partial structs.Test; this contains the test name, dir and categories.
func GetTestSummariesFromNames(names []string) []structs.Test {
	log.Printf("Reading tests from DB; want %d tests", len(names))

	nameListSQL := "'" + strings.Join(names, "', '") + "'"
	sql := "SELECT name, dir, categories from " + constants.RegisteredTestTable + " WHERE test.name in (" + nameListSQL + ")"
	log.Println("SQL :", sql)

	var tests []structs.Test

	log.Println("SQL :", sql)
	rows, err := DBConn.Query(sql)
	if err != nil {
		log.Criticalf("Error reading tests: %q", err)
	}
	for rows.Next() {
		var t structs.Test
		err = rows.Scan(&t.Name, &t.Dir, &t.Categories)
		if err != nil {
			log.Criticalf("Error reading row into struct: %q", err)
		}

		tests = append(tests, t)
	}

	log.Debugf("Retrieved %d tests in DB", len(tests))
	if len(names) != len(tests) {
		log.Errorf("Wanted summaries of %d tests, but only retrieved %d tests", len(names), len(tests))
	}

	return tests
}
