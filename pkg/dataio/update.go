package dataio

// import (
// 	"bytes"
// 	"encoding/csv"
// 	_ "encoding/json"
// 	"fmt"
// 	_ "html/template"
// 	_ "path/filepath"

// 	_ "database/sql"

// 	_ "github.com/gorilla/websocket"
// 	_ "github.com/lib/pq"

// 	// "io/ioutil"
// 	_ "encoding/csv"
// 	log "github.com/romana/rlog"
// 	_ "net/http"
// 	"os"
// 	"ted/pkg/constants"
// 	_ "ted/pkg/handler" // TODO enable
// 	"ted/pkg/help"
// 	"ted/pkg/structs"
// 	_ "ted/pkg/structs"
// 	"ted/pkg/ws"
// 	_ "time"

// log "github.com/romana/rlog"
// )

// TODO this is a placeholder file

// func WriteResultToStore(result structs.Result) {
// 	if help.IsLocal {
// 		WriteResultToCSV(result)
// 	} else {
// 		WriteResultToDB(result)
// 	}
// 	log.Println("Result written to store")
// 	SendReload(result) // after writing, reload the page so that it shows the new results
// 	log.Println("After SendReload")
// }

// func SendReload(result structs.Result) {
// 	log.Println("Will try to send result to WS")
// 	message := result.ToJSON()
// 	messageBytes := bytes.TrimSpace([]byte(message))
// 	ws.WSHub.Broadcast <- messageBytes

// 	log.Println("Result sent to WS: ", message)
// }

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

// func WriteResultToDB(result structs.Result) {
// 	log.Println("Writing result to DB")
// 	sql := constants.ResultsTableInsertSQL + fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s')", result.TestName, result.TestRunIdentifier, result.Category, result.Status, result.Timestamp, result.Message)
// 	log.Println("SQL :", sql)
// 	if _, err := DBConn.Exec(sql); err != nil {
// 		log.Criticalf("Error writing result to DB: %q", err)
// 	}
// 	// TODO

// }
