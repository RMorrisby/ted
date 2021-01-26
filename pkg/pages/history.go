package pages

import (
	"net/http"
	"ted/pkg/dataio"
	"ted/pkg/enums"
	"ted/pkg/help"
	"ted/pkg/structs"

	log "github.com/romana/rlog"
)

// Page showing the summary historical view of a test suite
func HistoryPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("suite")
	if name == "" {
		// If no suite has been specified, default to the suite for the latest result
		name = dataio.LatestSuite
	}

	suite := dataio.GetSuite(name)
	if suite == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Suite '" + name + "' is not registered in TED"))
		return
	}

	// ws.ServeWs(ws.WSHub, w, r)

	// err := Templates.ExecuteTemplate(w, "history.html", GetHistoryForSuite(name)) //execute the template and pass it the struct to fill in the gaps
	err := Templates.ExecuteTemplate(w, "history.html", name)

	if err != nil {
		log.Debug("template executing error: ", err)
	}
}

// REST endpoint to get the history data for a suite, to be displayed in the UI
func HistoryOfSuite(w http.ResponseWriter, r *http.Request) {

	help.LogNewAPICall("HistoryOfSuite")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("suite")
	if name == "" {
		// A suite name must be supplied
		s := "No suite name supplied to " + r.Method + " " + r.URL.RequestURI() + "; URL must be /historydata?suite=___"
		log.Error(s)
		http.Error(w, s, http.StatusBadRequest)
		return
	}

	history := GetHistoryForSuite(name)

	help.MarshalJSONAndWriteToResponse(history, w)
}

// REST endpoint to get the history data for a suite, to be displayed in the UI
// This deliberately only returns the history data for the last 5 test runs, so that they fit in the UI
// TODO fix the UI page so all data can be displayed, with the header & test-name columns always visible
func HistoryOfSuiteRecent(w http.ResponseWriter, r *http.Request) {

	help.LogNewAPICall("HistoryOfSuiteRecent")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("suite")
	if name == "" {
		// A suite name must be supplied
		s := "No suite name supplied to " + r.Method + " " + r.URL.RequestURI() + "; URL must be /historydatarecent?suite=___"
		log.Error(s)
		http.Error(w, s, http.StatusBadRequest)
		return
	}

	history := GetHistoryForSuite(name)
	count := len(history.TestRuns)
	maxCount := 6
	if count > maxCount {
		history.TestRuns = history.TestRuns[count-maxCount:]
		for testrun := range history.TestRunMap {
			if !help.Contains(history.TestRuns, testrun) {
				delete(history.TestRunMap, testrun)
			}
		}
	}

	help.MarshalJSONAndWriteToResponse(history, w)
}

// Get the test history for the given suite
// Contains all the info needed to populate the suite-history page
func GetHistoryForSuite(suiteName string) structs.HistorySuite {

	var history structs.HistorySuite

	// This list should already be sorted by the testrun (ascending)
	allResults := dataio.ReadAllResultsForSuite(suiteName)
	// log.Debug("ALL RESULTS ")
	// log.Debug(allResults)
	// Some validation, just to be sure
	for _, result := range allResults {
		if result.SuiteName != suiteName {
			log.Errorf("Retrieval of all results for suite %s included a result for suite %s :: %s", suiteName, result.SuiteName, result.ToJSON())
		}
	}

	// Parse the results to build the TestName::HistoryTestSummary map

	/*
		The matrix of results for each test run may be sparse (e.g. test are added to the suite in more recent test runs)
		Therefore we need to insert fake results (NOT RUN) in those gaps
		To identify gaps, we need both the list of tests for the suite, and the list of test runs for the suite
	*/

	// Separate the results out, collecting them under each test

	// Each list in this map may be of different length (because the results matrix may be sparse)
	tempSparseResultsMap := make(map[string][]structs.Result) // TestRun::[]Result

	// Separate the results out by their test run
	// log.Debug("NOW PARSING RESULTS")
	for _, result := range allResults {
		// log.Debug("")
		// log.Debugf("Result %s::%s | %s", result.TestName, result.TestRunIdentifier, tempSparseResultsMap[result.TestRunIdentifier])
		if _, ok := tempSparseResultsMap[result.TestRunIdentifier]; ok {
			tempSparseResultsMap[result.TestRunIdentifier] = append(tempSparseResultsMap[result.TestRunIdentifier], result)
		} else {
			tempSparseResultsMap[result.TestRunIdentifier] = []structs.Result{result}
		}
		// log.Debug(tempSparseResultsMap[result.TestRunIdentifier])
		// log.Debug("")
	}
	// log.Debugf("FIN")

	// Get the list of test names
	// This list will not be in the desired final order, so we will need to rebuild it later
	var tempTestNames []string // A set of test names
	for _, r := range allResults {
		if !help.Contains(tempTestNames, r.TestName) {
			tempTestNames = append(tempTestNames, r.TestName)
		}
	}

	// Get the summary of each test & store in history
	tests := dataio.GetTestSummariesFromNames(tempTestNames)
	var testNames = make([]string, len(tests))

	for i, t := range tests {
		testNames[i] = t.Name
	}

	// Each list in this map will be the same length (extra results will be inserted to eliminate sparseness)
	tempResultsMap := make(map[string][]structs.Result) // TestRun::[]Result // will have extra results inserted where necessary
	// Initialise the non-sparse results map
	for testRun := range tempSparseResultsMap {
		tempResultsMap[testRun] = []structs.Result{}
	}

	// We also need to use this list to sort the results for each test
	testRuns := dataio.GetAllTestRunsForSuite(suiteName)

	// This needs to be approriately sorted
	help.SortTestRuns(testRuns)

	// Fill out the non-sparse results map with fake results where necessary
	// Add each result in order of each testrun
	// log.Debug("NOW ADDING UNKNOWN RESULTS")
	// log.Debug("")
	// log.Debug(tempResultsMap)
	for _, testRun := range testRuns {
		// log.Debug("Test run : ", testRun)
		for _, testName := range testNames {
			// log.Debug("Testname : ", testName)
			// log.Debug("")
			knownResults := tempSparseResultsMap[testRun]
			// log.Debug("testrun         : ", testRun)
			// log.Debug("known results   : ", knownResults)
			// log.Debug("contains result : ", ContainsResultForTest(knownResults, testName))
			if ContainsResultForTest(knownResults, testName) {
				tempResultsMap[testRun] = append(tempResultsMap[testRun], GetResultForTestFromList(knownResults, testName))
			} else {
				var fakeResult structs.Result
				fakeResult.TestRunIdentifier = testRun
				fakeResult.SuiteName = suiteName
				fakeResult.TestName = testName
				fakeResult.Status = "UNKNOWN"
				fakeResult.TedStatus = "UNKNOWN"
				tempResultsMap[testRun] = append(tempResultsMap[testRun], fakeResult)
			}
		}
	}
	// Validate that each list of results is the same length, equal to the length of the set of testruns

	// for i, test := range tests {
	// 	log.Debugf("Test %d :: %s", i+1, test.Name)
	// }

	// for testRun, results := range tempResultsMap {
	// 	log.Debugf("Test run %s :: %d results", testRun, len(results))
	// 	log.Debug(results)
	// }

	testRunMap := make(map[string]structs.HistoryTestRun)

	for testRun, resultList := range tempResultsMap {
		var summary structs.HistoryTestRun
		summary.TestRunName = testRun
		summary.ResultList = resultList

		var total int // total number of executed tests (doesn't include tests that weren't run)
		var success int
		var fail int
		var notRun int
		var knownIssue int

		for _, result := range summary.ResultList {

			switch result.TedStatus {
			case "PASSED":
				success++
				total++
			case enums.Failed:
				fail++
				total++
			case enums.KnownIssue:
				knownIssue++
				fail++ // Known Issue tests will have failed, by definition
				total++
			case enums.NotRun:
				notRun++
			}
		}

		summary.TotalCount = total
		summary.SuccessCount = success
		summary.FailCount = fail
		summary.NotRunCount = notRun
		summary.KnownIssueCount = knownIssue

		testRunMap[testRun] = summary
	}

	history.SuiteName = suiteName
	history.TestRuns = testRuns
	history.TestRunMap = testRunMap
	history.Tests = tests

	lastTestRunName := history.TestRuns[len(history.TestRuns)-1]
	lastTestRunSummary := history.TestRunMap[lastTestRunName]

	history.TotalCount = lastTestRunSummary.TotalCount
	history.SuccessCount = lastTestRunSummary.SuccessCount
	history.FailCount = lastTestRunSummary.FailCount
	history.NotRunCount = lastTestRunSummary.NotRunCount
	history.KnownIssueCount = lastTestRunSummary.KnownIssueCount

	// message, _ := json.Marshal(history)
	// log.Debug("History JSON")
	// log.Debug(string(message))
	// log.Debug("")
	// log.Debug("")
	return history
}

// Contains asks whether the Result list contains a result for the supplied test
func ContainsResultForTest(results []structs.Result, testname string) bool {
	for _, result := range results {
		if result.TestName == testname {
			return true
		}
	}
	return false
}

// Return the result in the list with the matching test
func GetResultForTestFromList(results []structs.Result, testname string) structs.Result {
	for _, result := range results {
		if result.TestName == testname {
			return result
		}
	}
	// go won't let us return nil
	return structs.Result{}
}
