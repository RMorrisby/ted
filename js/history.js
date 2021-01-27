// ############# Websocket stuff
// TODO

// ############# END Websocket stuff

// TODO other stuff?

// function appendKnownIssueButtonsToLastResults() {
//   var table = document.getElementById("history-table-body");

//   var last = $(table).find(".child").find("td").last();

//   // TODO
// }

function scrollRightMax() {
  var scrollWidth = $("#history-table-body").scrollWidth;
  console.log(scrollWidth);
  $("#history-table-body").scrollLeft(scrollWidth);
}

// Empties the table, so that it can be reformed
function clearHistoryTable() {
  $("div#history-div").empty();
  // var div = $("div#history-div");
  // if (div.hasChildNodes()) {
  //   while (div.hasChildElements()) {
  //     div.removeChild(div.lastChild);
  //   }
  // }
}

// Gets the history info from TED and builds the history table
function getFullHistory() {
  var div = $("div#history-div");
  var suite = div.attr("suite");
  console.log(suite);

  clearHistoryTable();

  $.get("/historydata?suite=" + suite, function (data) {
    console.log("Received history data for suite " + suite);
    buildHistoryTable(data);
  });

  // Change the button so that it will now get the recent history
  var button = $("button#history-toggle");
  button.text("Get recent history for suite " + suite);
  button.attr("onclick", "getRecentHistory()");
}

// Gets the history info for the most recent X test runs from TED and builds the history table
function getRecentHistory() {
  var div = $("div#history-div");
  var suite = div.attr("suite");
  console.log(suite);

  clearHistoryTable();

  $.get("/historydatarecent?suite=" + suite, function (data) {
    console.log("Received history data for suite " + suite);
    buildHistoryTable(data);
  });

  // Change the button so that it will now get the full history
  var button = $("button#history-toggle");
  button.text("Get full history for suite " + suite);
  button.attr("onclick", "getFullHistory()");
}

// Build the History table HTML, using the supplied data (a HistorySuite struct)
function buildHistoryTable(data) {
  // SuiteName       string
  // TestRuns        []string                  // list of all of the test runs the suite has any results for
  // TestRunMap      map[string]HistoryTestRun // map of the test run name and its results (Result struct)
  // Tests           []Test                    // (Test Summary) list of all of the tests the suite has any results for
  // TotalCount      int                       // the total number of executed tests in the most recent test run
  // SuccessCount    int                       // the number of test successes in the most recent test run
  // FailCount       int                       // the number of test failures in the most recent test run
  // KnownIssueCount int                       // the number of tests with known issues in the most recent test run
  // NotRunCount     int                       // the number of tests not yet run in the most recent test run

  // TestRunName     string
  // ResultList      []Result
  // TotalCount      int // the total number of executed tests in this test run
  // SuccessCount    int // the number of test successes in this test run
  // FailCount       int // the number of test failures in this test run
  // KnownIssueCount int // the number of tests with known issues in this test run
  // NotRunCount     int // the number of tests not yet run in this test run

  var json = JSON.parse(data);
  console.log(`Received history for ${json.TestRuns.length} test runs`);

  // Insert all elements into the main div
  var body = document.getElementById("history-div");

  // Build the stats summary
  var suiteStatsDiv = document.createElement("div");
  suiteStatsDiv.id = "suite-stats-div";
  suiteStatsDiv.appendChild(document.createTextNode("Showing results for suite " + json.SuiteName));
  // TODO the rest of this div
  body.appendChild(suiteStatsDiv);

  // Build the table header
  var table = document.createElement("table");
  table.id = "history-table";
  var head = document.createElement("thead");
  var tr = document.createElement("tr");
  tr.id = "history-table-header";

  var th = document.createElement("th");
  th.className = "history-table-dir";
  th.id = "history-table-header-dir";
  th.appendChild(document.createTextNode("Dir"));
  tr.appendChild(th);

  var th = document.createElement("th");
  th.id = "history-table-header-test";
  th.appendChild(document.createTextNode("Test"));
  tr.appendChild(th);

  var th = document.createElement("th");
  th.id = "history-table-header-categories";
  th.appendChild(document.createTextNode("Categories"));
  tr.appendChild(th);

  // Add headers for each test run
  for (var i = 0; i < json.TestRuns.length; i++) {
    var testrun = json.TestRuns[i];

    var th = document.createElement("th");
    th.id = testrun;
    var a = document.createElement("a");
    a.appendChild(document.createTextNode(testrun));
    a.title = testrun;
    a.href = "data?testrun=" + downcaseAndUnderscore(testrun);
    th.appendChild(a);
    tr.appendChild(th);
  }
  // Header for the Known Issue buttons & input field
  var th = document.createElement("th");
  th.className = "history-table-header-known-issue";
  th.appendChild(document.createTextNode("Known Issue?"));
  tr.appendChild(th);

  head.appendChild(tr);

  // Build the table body
  var tbody = document.createElement("tbody");
  tbody.id = "history-table-body";

  for (var i = 0; i < json.Tests.length; i++) {
    var test = json.Tests[i];
    var testNameDown = downcaseAndUnderscore(test.Name);

    // Create the tr
    var tr = document.createElement("tr");
    tr.className = "history-table-row";
    tr.id = "history-table-row-" + testNameDown;

    var td = document.createElement("td");
    td.className = "history-table-dir";
    td.id = "history-table-dir-" + testNameDown;
    td.appendChild(document.createTextNode(test.Dir));
    tr.appendChild(td);

    var td = document.createElement("td");
    td.id = "history-table-test-" + testNameDown;
    td.appendChild(document.createTextNode(test.Name));
    tr.appendChild(td);

    var td = document.createElement("td");
    td.id = "history-table-categories-" + testNameDown;
    td.appendChild(document.createTextNode(test.Categories));
    tr.appendChild(td);

    var lastTestRun = json.TestRuns.slice(-1).pop();

    // Write each result for the test
    for (var k = 0; k < json.TestRuns.length; k++) {
      var testrun = json.TestRuns[k];
      var testrunHistory = json.TestRunMap[testrun];
      var result = testrunHistory.ResultList[i];

      var td = document.createElement("td");
      // Give the cell two statuses - the test status and the TED status
      // The TED status takes precedence for controlling the cell's formatting, with the test status as the backup
      td.classList.add("test-" + downcaseAndUnderscore(result.Status));
      td.classList.add("test-" + downcaseAndUnderscore(result.TedStatus));

      // If this is the last testrun, also give the cell a fixed class that can be used to get the cell (within the row)
      if (testrun == lastTestRun) {
        td.classList.add("tedstatus");
      }
      td.id = "history-table-" + downcaseAndUnderscore(result.TestRunIdentifier);
      td.appendChild(document.createTextNode(result.TedStatus));
      tr.appendChild(td);
    }

    addKnownIssueFieldsToTableRow(tr, test.Name, lastTestRun, test.KnownIssueDescription);
    tbody.appendChild(tr);
  }
  console.log(data); // TODO remove

  table.appendChild(head);
  table.appendChild(tbody);
  body.appendChild(table);
}

// On page load, adorn the table with whstever extra elements we need
// JS requires this function-wrapping
window.onload = function () {
  getRecentHistory();
  // scrollRightMax();
};
