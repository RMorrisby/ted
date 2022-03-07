package constants

const (
	LayoutDateISO = "2006-01-02"
	LayoutTimeISO = "15:04:05"

	ResultCSVFilename = "result.csv"

	// Suite Table
	SuiteTable                  = "suite"
	SuiteTableColumnDefinitions = "id serial primary key, name varchar(100) unique, description varchar(255), owner varchar(32), notes varchar(100)"
	SuiteTableCreateSQL         = "CREATE TABLE IF NOT EXISTS " + SuiteTable + " (" + SuiteTableColumnDefinitions + ")"
	SuiteTableInsertFullRowSQL  = "INSERT INTO " + SuiteTable + " (name, description, owner, notes) VALUES "

	SuiteTableSelectAllSQL  = "SELECT name, description, owner, notes from " + SuiteTable
	SuiteTableSelectNameSQL = "SELECT name from " + SuiteTable

	// Registered Test table
	RegisteredTestTable                  = "test"
	RegisteredTestTableColumnDefinitions = "id serial primary key, name varchar(100) unique, dir varchar(32), priority integer, categories varchar(255), description varchar(255), notes varchar(100), owner varchar(32), is_known_issue boolean default false, known_issue_description varchar(255)"
	RegisteredTestTableCreateSQL         = "CREATE TABLE IF NOT EXISTS " + RegisteredTestTable + " (" + RegisteredTestTableColumnDefinitions + ")"
	RegisteredTestTableInsertFullRowSQL  = "INSERT INTO " + RegisteredTestTable + " (name, dir, priority, categories, description, notes, owner, is_known_issue, known_issue_description) VALUES "

	RegisteredTestTableSelectAllSQL  = "SELECT name, dir, priority, categories, description, notes, owner, is_known_issue, known_issue_description from " + RegisteredTestTable
	RegisteredTestTableSelectNameSQL = "SELECT name from " + RegisteredTestTable

	// Result table
	ResultTable                  = "result"
	ResultTableColumnDefinitions = "id serial primary key, suite_id integer references " + SuiteTable + ", test_id integer references " + RegisteredTestTable + ", testrun varchar(32) not null, status varchar(32), start_time timestamp with time zone, end_time timestamp with time zone, ran_by varchar(32), message varchar(100), ted_status varchar(32), ted_notes varchar(255)"
	ResultTableCreateSQL         = "CREATE TABLE IF NOT EXISTS " + ResultTable + " (" + ResultTableColumnDefinitions + ")"
	ResultTableInsertFullRowSQL  = "INSERT INTO " + ResultTable + " (suite_id, test_id, testrun, status, start_time, end_time, ran_by, message, ted_status, ted_notes) VALUES "

	ResultTableInsertNotYetRunRowSQL = "INSERT INTO " + ResultTable + " (suite_id, test_id, testrun, ted_status, ted_notes) VALUES "

	// Reads all results from the DB, yielding the fields that the Result struct wants (i.e. test.name instead of result.test_id)
	ResultTableSelectAllSQL = "SELECT suite.name, test.name, result.testrun, result.status, COALESCE(to_char(result.start_time, 'YYYY-MM-DD HH24:MI:SS'), ''), COALESCE(to_char(result.end_time, 'YYYY-MM-DD HH24:MI:SS'), ''), result.ran_by, result.message, result.ted_status, result.ted_notes FROM " + ResultTable + " result LEFT JOIN " + SuiteTable + " suite ON result.suite_id = suite.id LEFT JOIN " + RegisteredTestTable + " test ON result.test_id = test.id ORDER BY suite.name ASC, result.testrun ASC, test.name ASC"

	// Reads all results from the DB, without an ORDER BY clause
	// You may append WHERE clauses, etc., to the end of this

	ResultTableSelectAllNoSortingSQL = "SELECT suite.name, test.name, result.testrun, result.status, COALESCE(to_char(result.start_time, 'YYYY-MM-DD HH24:MI:SS'), ''), COALESCE(to_char(result.end_time, 'YYYY-MM-DD HH24:MI:SS'), ''), result.ran_by, result.message, result.ted_status, result.ted_notes FROM " + ResultTable + " result LEFT JOIN " + SuiteTable + " suite ON result.suite_id = suite.id LEFT JOIN " + RegisteredTestTable + " test ON result.test_id = test.id"

	// Reads all results from the DB, yielding the fields that the ResultForUI struct wants (i.e. test.name instead of result.test_id)
	ResultTableSelectAllResultsForUISQL = "SELECT test.categories, test.dir, test.name, result.testrun, result.status, test.priority, COALESCE(to_char(result.start_time, 'YYYY-MM-DD HH24:MI:SS'), ''), COALESCE(to_char(result.end_time, 'YYYY-MM-DD HH24:MI:SS'), ''), result.ran_by, result.message, result.ted_status, result.ted_notes FROM " + ResultTable + " result LEFT JOIN " + RegisteredTestTable + " test ON result.test_id = test.id ORDER BY result.testrun ASC, test.name ASC"

	// Reads all results and returns the set of test runs
	ResultTableSelectDistinctTestRunNoSortingSQL = "SELECT DISTINCT ON (result.testrun) result.testrun FROM " + ResultTable + " result LEFT JOIN " + SuiteTable + " suite ON result.suite_id = suite.id"
)
