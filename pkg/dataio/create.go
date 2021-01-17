package dataio

import (
	"bytes"
	"fmt"
	"ted/pkg/constants"
	_ "ted/pkg/handler" // TODO enable
	_ "ted/pkg/help"
	"ted/pkg/structs"
	"ted/pkg/ws"

	log "github.com/romana/rlog"
)

func WriteResultToStore(result structs.Result) {
	// if help.IsLocal {
	// 	WriteResultToCSV(result)
	// } else {
	resultForUI := WriteFullResultToDB(result)
	// }
	log.Println("Result written to store")
	SendReload(resultForUI) // after writing, reload the page so that it shows the new results
	log.Println("After SendReload")
}

func SendReload(resultForUI structs.ResultForUI) {
	log.Println("Will try to send resultForUI to WS")
	message := resultForUI.ToJSON()
	messageBytes := bytes.TrimSpace([]byte(message))
	ws.WSHub.Broadcast <- messageBytes

	log.Println("Result sent to WS: ", message)
}

// func WriteResultToCSV(result structs.Result) {
// 	log.Println("Will now write result to file :", result)
// 	// TODO use PSV instead of CSV
// 	// TODO don't write duplicates?
// 	f, err := os.OpenFile(constants.ResultCSVFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer f.Close()

// 	writer := csv.NewWriter(f)
// 	defer writer.Flush()

// 	resultArray := result.ToA()

// 	err = writer.Write(resultArray)
// 	help.CheckError("Cannot write to file", err)

// 	log.Println("Wrote result to file")
// }

func WriteFullResultToDB(result structs.Result) (resultForUI structs.ResultForUI) {

	suiteID := fmt.Sprintf("(SELECT id from suite where suite.name = '%s')", result.SuiteName)
	testID := fmt.Sprintf("(SELECT id from test where test.name = '%s')", result.TestName)

	// TODO Maybe try something like this?
	// "INSERT INTO " + ResultTable + " ((select suite.id from suite where suite.name is " + result.SuiteName + ")),
	// (select test.id from test where test.name is " + result.Name + ")),
	// testrun, status, start_time, end_time, ran_by, message, ted_status, ted_notes) VALUES "

	log.Println("Writing result to DB")
	// (suite_id, test_id, testrun, status, start_time, end_time, ran_by, message, ted_status, ted_notes)
	sql := constants.ResultTableInsertFullRowSQL + fmt.Sprintf("(%s, %s, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')", suiteID, testID, result.TestRunIdentifier, result.Status, result.StartTimestamp, result.EndTimestamp, result.RanBy, result.Message, result.Status, "")
	log.Println("SQL :", sql)
	if _, err := DBConn.Exec(sql); err != nil {
		log.Criticalf("Error writing result to DB: %q", err)
	}

	// Now gather the info we need for the ResultForUI object
	// Get the test
	test := GetTest(result.TestName)

	resultForUI.TestName = result.TestName
	resultForUI.TestRunIdentifier = result.TestRunIdentifier
	resultForUI.Status = result.Status
	resultForUI.StartTimestamp = result.StartTimestamp
	resultForUI.EndTimestamp = result.EndTimestamp
	resultForUI.RanBy = result.RanBy
	resultForUI.Message = result.Message
	resultForUI.TedStatus = result.TedStatus
	resultForUI.TedNotes = result.TedNotes

	resultForUI.Categories = test.Categories
	resultForUI.Dir = test.Dir
	resultForUI.Priority = test.Priority
	return
}

// Write the suite to the DB, if the DB does not already contain a suite of that name
func WriteSuiteToDBIfNew(suite structs.Suite) {

	if !SuiteExists(suite.Name) {

		log.Println("Writing suite to DB")
		sql := constants.SuiteTableInsertFullRowSQL + fmt.Sprintf("('%s', '%s', '%s', '%s')", suite.Name, suite.Description, suite.Owner, suite.Notes)
		log.Println("SQL :", sql)
		if _, err := DBConn.Exec(sql); err != nil {
			log.Criticalf("Error writing result to DB: %q", err)
		}
	} else {
		log.Printf("Suite %s already exists", suite.Name)
	}
}

// Write the test to the DB, if the DB does not already contain a test of that name
func WriteTestToDBIfNew(test structs.Test) {

	if !TestExists(test.Name) {

		log.Println("Writing test to DB")
		// (name, dir, priority, categories, description, notes, owner, is_known_issue, known_issue_description) VALUES "
		sql := constants.RegisteredTestTableInsertFullRowSQL + fmt.Sprintf("('%s', '%s', '%d', '%s', '%s', '%s', '%s', '%t', '%s')",
			test.Name, test.Dir, test.Priority, test.Categories, test.Description, test.Notes, test.Owner, test.IsKnownIssue, test.KnownIssueDescription)
		log.Println("SQL :", sql)
		if _, err := DBConn.Exec(sql); err != nil {
			log.Criticalf("Error writing result to DB: %q", err)
		}
	} else {
		log.Printf("Test %s already exists", test.Name)
	}
}
