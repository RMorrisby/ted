package dataio

import (
	"fmt"
	"ted/pkg/constants"
	"ted/pkg/enums"
	"ted/pkg/help"
	"ted/pkg/structs"

	log "github.com/romana/rlog"
)

// Update the Known Issue fields for the given test
func WriteTestKnownIssueUpdate(update structs.KnownIssueUpdate) {

	log.Println("Updating test in DB")
	sql := fmt.Sprintf("UPDATE %s SET is_known_issue = %t, known_issue_description = '%s' WHERE name = '%s'", constants.RegisteredTestTable, update.IsKnownIssue, update.KnownIssueDescription, update.TestName)
	log.Println("SQL :", sql)
	if _, err := DBConn.Exec(sql); err != nil {
		log.Criticalf("Error writing result to DB: %q", err)
	}
}

// Update the Known Issue fields for the given result
func WriteResultKnownIssueUpdate(update structs.KnownIssueUpdate) {

	log.Println("Updating result in DB")

	var tedStatus string
	var tedNotes string
	if update.IsKnownIssue {
		tedStatus = enums.KnownIssue
		tedNotes = update.KnownIssueDescription
	} else {
		test := ReadResult(update.TestName, update.TestRun)
		tedStatus = test.Status // reset the TedStatus to the Status of the test // Might be PASSED, might be FAILED
		tedNotes = ""
	}

	testID := fmt.Sprintf("(SELECT id FROM test WHERE test.name = '%s')", update.TestName)

	sql := fmt.Sprintf("UPDATE %s SET ted_status = '%s', ted_notes = '%s' WHERE test_id = %s AND testrun = '%s'", constants.ResultTable, tedStatus, tedNotes, testID, update.TestRun)
	log.Println("SQL :", sql)
	if _, err := DBConn.Exec(sql); err != nil {
		log.Criticalf("Error writing result to DB: %q", err)
	}
}

// Overwrite the existing result with the given result
func WriteResultUpdate(update structs.Result, existing *structs.Result) structs.ResultForUI {

	log.Debug("Result update :", update.ToJSON())

	// suiteID := fmt.Sprintf("(SELECT id from suite where suite.name = '%s')", update.SuiteName)
	testID := fmt.Sprintf("(SELECT id FROM test WHERE test.name = '%s')", update.TestName)

	log.Println("Updating result in DB")

	// (suite_id, test_id, testrun, status, start_time, end_time, ran_by, message, ted_status, ted_notes)

	var tedStatus string
	var tedNotes string

	tedNotes = update.TedNotes // TedNotes is == Known Issue // TODO do anything more with this field?
	tedStatus = update.TedStatus
	// Failed -> passed
	if existing.Status == string(enums.Failed) && update.Status == string(enums.Passed) {
		tedStatus = string(enums.PassedOnRerun)
	}
	// Passed -> Failed // why was this rerun? On a whim?
	if existing.Status == string(enums.Passed) && update.Status == string(enums.Failed) {
		tedStatus = string(enums.Intermittent)
	}
	log.Debugf("TED Status : %s; TED Notes : %s", tedStatus, tedNotes)

	// sql := fmt.Sprintf("UPDATE %s result SET result.status = '%s', result.start_time = '%s', result.end_time = '%s', result.ran_by = '%s', result.message = '%s', result.ted_status = '%s', result.ted_notes = '%s' WHERE result.test_id = %s AND result.testrun = '%s'", constants.ResultTable, update.Status, update.StartTimestamp, update.EndTimestamp, update.RanBy, update.Message, tedStatus, tedNotes, testID, update.TestRunIdentifier)

	sql := fmt.Sprintf("UPDATE %s SET status = '%s', start_time = '%s', end_time = '%s', ran_by = '%s', message = '%s', ted_status = '%s', ted_notes = '%s' WHERE test_id = %s AND testrun = '%s'", constants.ResultTable, update.Status, update.StartTimestamp, update.EndTimestamp, update.RanBy, update.Message, tedStatus, tedNotes, testID, update.TestRunIdentifier)

	log.Println("SQL :", sql)
	if _, err := DBConn.Exec(sql); err != nil {
		log.Criticalf("Error updating result in DB: %q", err)
	}
	// Now gather the info we need for the ResultForUI object
	// Get the test
	test := GetTest(update.TestName)
	return help.FormResultForUI(update, test)
}
