package constants

const (
	LayoutDateISO = "2006-01-02"
	LayoutTimeISO = "15:04:05"

	ResultCSVFilename = "result.csv"

	ResultsTable                  = "results"
	ResultsTableColumnDefinitions = "id serial, name varchar(100), testrun varchar(32), category varchar(32), status varchar(32), endtime timestamp with time zone, message varchar(100)"
	ResultsTableCreateSQL         = "CREATE TABLE IF NOT EXISTS " + ResultsTable + " (" + ResultsTableColumnDefinitions + ")"
	ResultsTableInsertSQL         = "INSERT INTO " + ResultsTable + "(name, testrun, category, status, endtime, message) VALUES"

	ResultsTableSelectSQL         = "SELECT name, testrun, category, status, endtime, message from " + ResultsTable
)
