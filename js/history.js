// ############# Websocket stuff
// TODO

// ############# END Websocket stuff

// TODO other stuff?

function appendKnownIssueButtonsToLastResults() {
  var table = document.getElementById("history-table-body");

  var last = $(table).find(".child").find("td").last();

  // TODO
}

function scrollRightMax() {
  var scrollWidth = $("#history-table-body").scrollWidth;
  console.log(scrollWidth);
  $("#history-table-body").scrollLeft(scrollWidth);
}

// Gets the history info from TED and builds the history table
function buildHistoryTable() {
  var div = $("div#history-div");
  var suite = div.attr("suite");
  console.log(suite);

  $.get("/historydata?suite=" + suite, function (data) {
    console.log("Received all history data for suite " + suite);
    var json = JSON.parse(data);

    console.log(`Received history for ${json.TestRuns.length} test runs`);

    // SuiteName       string
    // TestRuns        []string                  // list of all of the test runs the suite has any results for
    // TestRunMap      map[string]HistoryTestRun // map of the test run name and its results (Result struct)
    // Tests           []Test                    // list of all of the tests the suite has any results for
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

    // Insert all elements into the main div
    var body = document.getElementById("history-div");

    // Build the stats summary
    var suiteStatsDiv = document.createElement("div");
    suiteStatsDiv.id = "suite-stats-div";
    suiteStatsDiv.appendChild(document.createTextNode("Showing results for suite " + suite));
    // TODO the rest of this div
    body.appendChild(suiteStatsDiv);

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

      // Create the sticky-div
      var stickyDiv = document.createElement("div");
      stickyDiv.className = "history-table-sticky-left";

      var td = document.createElement("td");
      td.className = "history-table-dir";
      td.id = "history-table-dir-" + testNameDown;
      td.appendChild(document.createTextNode(test.Dir));
      stickyDiv.appendChild(td);

      var td = document.createElement("td");
      td.id = "history-table-test-" + testNameDown;
      td.appendChild(document.createTextNode(test.Name));
      stickyDiv.appendChild(td);

      var td = document.createElement("td");
      td.id = "history-table-categories-" + testNameDown;
      td.appendChild(document.createTextNode(test.Categories));
      stickyDiv.appendChild(td);

      tr.appendChild(stickyDiv);

      // Write each result for the test
      for (var k = 0; k < json.TestRuns.length; k++) {
        var testrun = json.TestRuns[k];
        var testrunHistory = json.TestRunMap[testrun];
        var result = testrunHistory.ResultList[i];

        var td = document.createElement("td");
        td.className = "test-" + downcaseAndUnderscore(result.TedStatus);
        td.id = "history-table-" + downcaseAndUnderscore(result.TestRunIdentifier);
        td.appendChild(document.createTextNode(result.TedStatus));
        tr.appendChild(td);
      }

      // Button to clear the Known Issue value
      var buttonClear = document.createElement("button");
      buttonClear.className = "known-issue-clear";
      buttonClear.id = "history-table-button-known-issue-clear-" + testNameDown;
      buttonClear.appendChild(document.createTextNode("N"));
      buttonClear.setAttribute("test", test.Name);

      // Button to set the Known Issue value
      var buttonSet = document.createElement("button");
      buttonSet.className = "known-issue-set";
      buttonSet.id = "history-table-button-known-issue-set-" + testNameDown;
      buttonSet.appendChild(document.createTextNode("Y"));
      buttonSet.setAttribute("test", test.Name);

      // Input field for the Known Issue value
      var input = document.createElement("input");
      input.className = "known-issue-input";
      input.id = "history-table-input-known-issue-" + testNameDown;
      input.appendChild(document.createTextNode(test.KnownIssueDescription));
      input.setAttribute("test", test.Name);

      var td = document.createElement("td");
      td.appendChild(buttonClear);
      td.appendChild(buttonSet);
      td.appendChild(input);
      tr.appendChild(td);
      tbody.appendChild(tr);
    }
    // console.log(data);

    // Build the table header
    var table = document.createElement("table");
    table.id = "history-table";
    var head = document.createElement("thead");
    var tr = document.createElement("tr");
    tr.id = "history-table-header";

    var stickyHeaderDiv = document.createElement("div");
    stickyHeaderDiv.className = "history-table-sticky-left";
    var th = document.createElement("th");
    th.className = "history-table-dir";
    th.id = "history-table-header-dir";
    th.appendChild(document.createTextNode("Dir"));
    stickyHeaderDiv.appendChild(th);

    var th = document.createElement("th");
    th.id = "history-table-header-test";
    th.appendChild(document.createTextNode("Test"));
    stickyHeaderDiv.appendChild(th);

    var th = document.createElement("th");
    th.id = "history-table-header-categories";
    th.appendChild(document.createTextNode("Categories"));
    stickyHeaderDiv.appendChild(th);

    tr.appendChild(stickyHeaderDiv);

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
    // Headers for the Known Issue buttons & input field
    // var th = document.createElement("th");
    // th.className = "history-table-header-known-issue-clear";
    // tr.appendChild(th);
    // var th = document.createElement("th");
    // th.className = "history-table-header-known-issue-set";
    // tr.appendChild(th);
    var th = document.createElement("th");
    th.className = "history-table-header-known-issue-input";
    tr.appendChild(th);

    head.appendChild(tr);
    table.appendChild(head);
    table.appendChild(tbody);
    body.appendChild(table);
  });
}

// On page load, adorn the table with whstever extra elements we need
// JS requires this function-wrapping
window.onload = function () {
  buildHistoryTable();
  scrollRightMax();
};
