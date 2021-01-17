package dataio

import (
	"fmt"

	"ted/pkg/constants"
	_ "ted/pkg/handler" // TODO enable

	log "github.com/romana/rlog"
)

func DeleteAllResults() (success bool, err error) {
	// if help.IsLocal {
	// 	success, err = DeleteAllResultsCSV()
	// } else {
	success, err = DeleteAllResultsDB()
	// }
	return
}

// func DeleteAllResultsCSV() (bool, error) {

// 	log.Println("Will now delete results from file :", constants.ResultCSVFilename)

// 	// f, err := os.OpenFile(constants.ResultCSVFilename, os.O_TRUNC, perm)
// 	err := os.Truncate(constants.ResultCSVFilename, 0)
// 	if err != nil {
// 		return false, fmt.Errorf("could not open file %q for truncation: %v", constants.ResultCSVFilename, err)
// 	}

// 	InitResultsCSV(true)

// 	// if err = f.Close(); err != nil {
// 	// 	return fmt.Errorf("could not close file handler for %q after truncation: %v", constants.ResultCSVFilename, err)
// 	// }
// 	return true, nil
// }

func DeleteAllResultsDB() (bool, error) {

	log.Println("Will now delete all results from DB")

	sql := fmt.Sprintf("DELETE FROM %s", constants.ResultTable)
	log.Println("SQL :", sql)
	r, err := DBConn.Exec(sql)
	if err != nil {
		log.Criticalf("Error deleting all results: %q", err)
		return false, err
	}

	numDeleted, err := r.RowsAffected()
	if err != nil {
		log.Criticalf("Error deleting all results: %q", err)
	}

	log.Printf("Deleted %d results from the DB", numDeleted)

	return true, nil
}


func DeleteAllTests() (success bool, err error) {
	
	log.Println("Will now delete all tests from DB")
	// By definition, a result cannot exist without a test
	// Therefore we must delete all results before deleting all tests
	DeleteAllResultsDB()

	// Now delete all tests
	sql := fmt.Sprintf("DELETE FROM %s", constants.RegisteredTestTable)
	log.Println("SQL :", sql)
	r, err := DBConn.Exec(sql)
	if err != nil {
		log.Criticalf("Error deleting all tests: %q", err)
		return false, err
	}

	numDeleted, err := r.RowsAffected()
	if err != nil {
		log.Criticalf("Error deleting all tests: %q", err)
	}

	log.Printf("Deleted %d tests from the DB", numDeleted)

	return true, nil
}


func DeleteAllSuites() (success bool, err error) {
	
	log.Println("Will now delete all suites from DB")
	// By definition, a result cannot exist without a suite
	// Therefore we must delete all results before deleting all suites
	DeleteAllResultsDB()

	// Now delete all suites
	sql := fmt.Sprintf("DELETE FROM %s", constants.SuiteTable)
	log.Println("SQL :", sql)
	r, err := DBConn.Exec(sql)
	if err != nil {
		log.Criticalf("Error deleting all suites: %q", err)
		return false, err
	}

	numDeleted, err := r.RowsAffected()
	if err != nil {
		log.Criticalf("Error deleting all suites: %q", err)
	}

	log.Printf("Deleted %d suites from the DB", numDeleted)

	return true, nil
}