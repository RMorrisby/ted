package dataio

import (
	"fmt"

	"ted/pkg/constants"

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

func DeleteAllResultsForTest(testName string) (bool, error) {

	log.Println("Will now delete all results for test", testName, "from DB")

	sql := fmt.Sprintf("DELETE FROM %s r WHERE r.test_id = (SELECT id FROM %s t where t.name = '%s')", constants.ResultTable, constants.RegisteredTestTable, testName)
	log.Println("SQL :", sql)
	r, err := DBConn.Exec(sql)
	if err != nil {
		log.Criticalf("Error deleting results for test %s: %q", testName, err)
		return false, err
	}

	numDeleted, err := r.RowsAffected()
	if err != nil {
		log.Criticalf("Error deleting results for test %s: %q", testName, err)
	}

	log.Printf("Deleted %d results from the DB", numDeleted)

	return true, nil
}

func DeleteAllResultsForSuite(suiteName string) (bool, error) {

	log.Println("Will now delete all results for suite", suiteName, "from DB")

	sql := fmt.Sprintf("DELETE FROM %s r WHERE r.suite_id = (SELECT id FROM %s s where s.name = '%s')", constants.ResultTable, constants.SuiteTable, suiteName)
	log.Println("SQL :", sql)
	r, err := DBConn.Exec(sql)
	if err != nil {
		log.Criticalf("Error deleting results for suite %s: %q", suiteName, err)
		return false, err
	}

	numDeleted, err := r.RowsAffected()
	if err != nil {
		log.Criticalf("Error deleting results for suite %s: %q", suiteName, err)
	}

	log.Printf("Deleted %d results from the DB", numDeleted)

	return true, nil
}

func DeleteTestRun(testRunName string) (bool, error) {

	log.Println("Will now delete test run", testRunName, "from DB")

	sql := fmt.Sprintf("DELETE FROM %s r WHERE r.testrun = '%s'", constants.ResultTable, testRunName)
	log.Println("SQL :", sql)
	r, err := DBConn.Exec(sql)
	if err != nil {
		log.Criticalf("Error deleting test run %s: %q", testRunName, err)
		return false, err
	}

	numDeleted, err := r.RowsAffected()
	if err != nil {
		log.Criticalf("Error deleting test run %s: %q", testRunName, err)
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

func DeleteAllStatuses() (success bool, err error) {

	log.Println("Will now delete all statuses from DB")

	sql := fmt.Sprintf("DELETE FROM %s", constants.StatusTable)
	log.Println("SQL :", sql)
	r, err := DBConn.Exec(sql)
	if err != nil {
		log.Criticalf("Error deleting all statuses: %q", err)
		return false, err
	}

	numDeleted, err := r.RowsAffected()
	if err != nil {
		log.Criticalf("Error deleting all statuses: %q", err)
	}

	log.Printf("Deleted %d statuses from the DB", numDeleted)

	return true, nil
}

// Delete a specific test (and its results)
func DeleteTest(name string) (success bool, err error) {

	log.Println("Will now delete test", name, "from DB")
	// By definition, a result cannot exist without a test
	// Therefore we must delete all results relating to the test before deleting the test
	DeleteAllResultsForTest(name)

	// Now delete the suite
	sql := fmt.Sprintf("DELETE FROM %s WHERE name = '%s'", constants.RegisteredTestTable, name)
	log.Println("SQL :", sql)
	r, err := DBConn.Exec(sql)
	if err != nil {
		log.Criticalf("Error deleting test %s: %q", name, err)
		return false, err
	}

	numDeleted, err := r.RowsAffected()
	if numDeleted > 1 {
		log.Criticalf("%d tests deleted; should only have deleted one test", numDeleted)
	}
	if err != nil || numDeleted == 0 {
		log.Criticalf("Error deleting test %s: %q", name, err)
	}

	log.Printf("Deleted test %s from the DB", name)

	return true, nil
}

// Delete a specific suite (and its results)
func DeleteSuite(name string) (success bool, err error) {

	log.Println("Will now delete suite", name, "from DB")
	// By definition, a result cannot exist without a suite
	// Therefore we must delete all results relating to the suite before deleting the suite
	DeleteAllResultsForSuite(name)

	// Now delete the suite
	sql := fmt.Sprintf("DELETE FROM %s WHERE name = '%s'", constants.SuiteTable, name)
	log.Println("SQL :", sql)
	r, err := DBConn.Exec(sql)
	if err != nil {
		log.Criticalf("Error deleting suite %s: %q", name, err)
		return false, err
	}

	numDeleted, err := r.RowsAffected()
	if numDeleted > 1 {
		log.Criticalf("%d tests deleted; should only have deleted one test", numDeleted)
	}
	if err != nil || numDeleted == 0 {
		log.Criticalf("Error deleting suite %s: %q", name, err)
	}

	log.Printf("Deleted suite %s from the DB", name)

	return true, nil
}
