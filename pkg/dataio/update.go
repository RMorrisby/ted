package dataio

import (
	"fmt"
	"ted/pkg/constants"
	"ted/pkg/structs"
	"ted/pkg/enums"

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
	sql := fmt.Sprintf("UPDATE %s result LEFT JOIN %s test ON result.test_id = test.id SET result.ted_status = '%s', result.ted_notes = '%s' WHERE test.name = '%s' AND result.TestRunIdentifier = '%s'", constants.ResultTable, constants.RegisteredTestTable, tedStatus, tedNotes, update.TestName, update.TestRun)
	log.Println("SQL :", sql)
	if _, err := DBConn.Exec(sql); err != nil {
		log.Criticalf("Error writing result to DB: %q", err)
	}
}
