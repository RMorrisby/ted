package dataio

import (
	"bytes"
	_ "encoding/csv"
	_ "encoding/json"
	"fmt"
	_ "html/template"
	_ "path/filepath"

	_ "database/sql"

	_ "github.com/gorilla/websocket"
	_ "github.com/lib/pq"

	// "io/ioutil"
	_ "encoding/csv"
	"log"
	_ "net/http"
	_ "os"
	"ted/pkg/constants"
	_ "ted/pkg/handler" // TODO enable
	_ "ted/pkg/help"
	"ted/pkg/structs"
	_ "ted/pkg/structs"
	"ted/pkg/ws"
	_ "time"
)

func WriteResultToStore(result structs.Result) {
	// if help.IsLocal {
	// 	WriteResultToCSV(result)
	// } else {
	WriteFullResultToDB(result)
	// }
	log.Println("Result written to store")
	SendReload(result) // after writing, reload the page so that it shows the new results
	log.Println("After SendReload")
}

func SendReload(result structs.Result) {
	log.Println("Will try to send result to WS")
	message := result.ToJSON()
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

func WriteFullResultToDB(result structs.Result) {

	suiteID := "" // TODO use result.SuiteName
	testID := ""  // TODO use result.Name

	// TODO Maybe try something like this?
	// "INSERT INTO " + ResultTable + " ((select suite.id from suite where suite.name is " + result.SuiteName + ")),
	// (select test.id from test where test.name is " + result.Name + ")),
	// testrun, status, start_time, end_time, ran_by, message, ted_status, ted_notes) VALUES "

	log.Println("Writing result to DB")
	// (suite_id, test_id, testrun, status, start_time, end_time, ran_by, message, ted_status, ted_notes)
	sql := constants.ResultTableInsertFullRowSQL + fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s')", suiteID, testID, result.TestRunIdentifier, result.Status, result.StartTimestamp, result.EndTimestamp, result.RanBy, result.Message, result.Status, "")
	log.Println("SQL :", sql)
	if _, err := DBConn.Exec(sql); err != nil {
		log.Fatalf("Error writing result to DB: %q", err)
	}
}
