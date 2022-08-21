// JS file for common functions

// Returns the string with all whitespace replaced with underscores, and with all letters in lowercase
// TODO
function downcaseAndUnderscore(s) {
  // return s.replace(/\W+/g, "_").toLowerCase();
  // return s.replace(/\s+/g, "_").toLowerCase();
  return s.replace(/[^a-zA-Z0-9.-_=+]+/g, "_").toLowerCase();
}

function stringToBoolean(string) {
  switch (string.toLowerCase().trim()) {
    case "true":
    case "yes":
    case "1":
      return true;
    case "false":
    case "no":
    case "0":
    case null:
      return false;
    default:
      return Boolean(string);
  }
}

// Turns a timestamp in T-Z format (2021-01-17T19:41:00.000Z) into a more readable one (2021-01-17 19:41)
function makeTimestampHumanReadable(t) {
  if (!t) {
    return "";
  }
  console.log(t);
  // toISOString should yield a date in this format : 2021-01-17T19:41:00.000Z
  // We want 2021-01-17 19:41
  // TODO This needs to handle non-GMT timestamps properly - we're not displaying the timezone, so to
  // the user it looks like a local time
  // The page also needs to warn / declare this
  // Incredibly, JS doesn't have any handling for format-strings. So we have to brute-force this somewhat.
  var date = new Date(t).toISOString().replace(/(T|Z)/g, " ").slice(0, 16);
  return date;
}

// The statuses in the DB are all uppercase (e.g. FAILED, PASSED ON RERUN). These take up more space on the
// page than capitalised-case
function makeStatusesMoreReadable(status) {
  return status.charAt(0).toUpperCase() + status.slice(1).toLowerCase();
}

// To the table row e, append buttons to set the test result as Passed or Failed
// Takes : e : element (the tr)
//         testName : string (the test name)
//         lastTestRun : string (the value of the last test run)
function addUpdateTestStatusFieldsToTableRow(e, testName, lastTestRun) {
  var testNameDown = downcaseAndUnderscore(testName);

  // Button to set the test to Passed
  var buttonPassed = document.createElement("button");
  buttonPassed.className = "test-status-set-to-passed";
  buttonPassed.id = "history-table-button-test-status-set-to-passed-" + testNameDown;
  buttonPassed.appendChild(document.createTextNode("Pass"));
  buttonPassed.setAttribute("test", testName);
  buttonPassed.setAttribute("testrun", lastTestRun);
  buttonPassed.setAttribute("outcome", "PASSED");
  $(buttonPassed).on("click", function () {
    sendStatusUpdateForTestAndReloadPage(this);
  });
  // Button to set the test to Passed
  var buttonFailed = document.createElement("button");
  buttonFailed.className = "test-status-set-to-failed";
  buttonFailed.id = "history-table-button-test-status-set-to-failed-" + testNameDown;
  buttonFailed.appendChild(document.createTextNode("Fail"));
  buttonFailed.setAttribute("test", testName);
  buttonFailed.setAttribute("testrun", lastTestRun);
  buttonFailed.setAttribute("outcome", "FAILED");
  $(buttonFailed).on("click", function () {
    sendStatusUpdateForTestAndReloadPage(this);
  });

  var td = document.createElement("td");
  // TODO hide each button if the result is already that status (don't show 'Pass' when the test has passed)
  td.appendChild(buttonPassed);
  td.appendChild(buttonFailed);
  e.appendChild(td);
}

// Send to TED the updated status for the given test result
// And reloads the page
// TODO all we want to change is the class & value of the Status & TEDStatus clss for the current row, not reloading the whole page
function sendStatusUpdateForTestAndReloadPage(button) {
  var testName = button.getAttribute("test");
  var lastTestRun = button.getAttribute("testrun");
  var outcome = button.getAttribute("outcome");
  console.log("In sendSU; 1 2 3 :: " + testName + " :: " + lastTestRun + " :: " + outcome);

  $.ajax({
    url: "/result",
    method: "PATCH",
    contentType: "application/json; charset=utf-8",
    dataType: "json",

    data: JSON.stringify({
      TestName: testName,
      TestRunIdentifier: lastTestRun,
      Status: outcome,
    }),

    statusCode: {
      200: function (xhr) {
        reloadPage();
      },
    },
  });
}

// To the table row e, append the Known Issue fields
// Takes : e : element (the tr)
//         testName : string (the test name)
//         lastTestRun : string (the value of the last test run)
//         knownIssueDesc : string (the Known Issue for this test)
// function addKnownIssueFieldsToTableRow(e, test, lastTestRun) {
function addKnownIssueFieldsToTableRow(e, testName, lastTestRun, knownIssueDesc) {
  var testNameDown = downcaseAndUnderscore(testName);

  // Button to set the test to Passed
  var buttonPassed = document.createElement("button");
  buttonPassed.className = "test-status-set-to-passed";
  buttonPassed.id = "history-table-button-test-status-set-to-passed-" + testNameDown;
  buttonPassed.appendChild(document.createTextNode("Pass"));
  buttonPassed.setAttribute("test", testName);
  buttonPassed.setAttribute("testrun", lastTestRun);
  buttonPassed.setAttribute("outcome", "PASSED");
  $(buttonPassed).on("click", function () {
    sendStatusUpdateForTest(this);
  });
  // Button to set the test to Passed
  var buttonFailed = document.createElement("button");
  buttonFailed.className = "test-status-set-to-failed";
  buttonFailed.id = "history-table-button-test-status-set-to-failed-" + testNameDown;
  buttonFailed.appendChild(document.createTextNode("Fail"));
  buttonFailed.setAttribute("test", testName);
  buttonFailed.setAttribute("testrun", lastTestRun);
  buttonFailed.setAttribute("outcome", "FAILED");
  $(buttonFailed).on("click", function () {
    sendStatusUpdateForTest(this);
  });

  // Button to clear the Known Issue value
  var buttonClear = document.createElement("button");
  buttonClear.className = "known-issue-clear";
  buttonClear.id = "history-table-button-known-issue-clear-" + testNameDown;
  buttonClear.appendChild(document.createTextNode("N"));
  buttonClear.setAttribute("test", testName);
  buttonClear.setAttribute("testrun", lastTestRun);
  buttonClear.setAttribute("is-known-issue", false);
  $(buttonClear).on("click", function () {
    sendKnownIssueForTest(this);
  });

  // Button to set the Known Issue value
  var buttonSet = document.createElement("button");
  buttonSet.className = "known-issue-set";
  buttonSet.id = "history-table-button-known-issue-set-" + testNameDown;
  buttonSet.appendChild(document.createTextNode("Y"));
  buttonSet.setAttribute("test", testName);
  buttonSet.setAttribute("testrun", lastTestRun);
  buttonSet.setAttribute("is-known-issue", true);
  $(buttonSet).on("click", function () {
    sendKnownIssueForTest(this);
  });

  // Input field for the Known Issue value
  var input = document.createElement("input");
  input.className = "known-issue-input";
  input.id = "history-table-input-known-issue-" + testNameDown;
  input.value = knownIssueDesc;
  input.setAttribute("test", testName);

  var td = document.createElement("td");
  // td.appendChild(buttonPassed);
  // td.appendChild(buttonFailed);
  td.appendChild(buttonClear);
  td.appendChild(buttonSet);
  td.appendChild(input);
  e.appendChild(td);
}

// Send to TED the desired Known Issue value for the test, either setting it or clearing it
function sendKnownIssueForTest(button) {
  var testName = button.getAttribute("test");
  var testNameDown = downcaseAndUnderscore(testName);
  var lastTestRun = button.getAttribute("testrun");
  console.log("In setKI; 1 2 3 :: " + testNameDown + " :: " + testName + " :: " + lastTestRun);
  var desc = $("input#history-table-input-known-issue-" + testNameDown).val();
  console.log("Desc : " + desc);
  var isKnownIssue = stringToBoolean(button.getAttribute("is-known-issue"));
  // If we are clearing the Known Issue, set desc to ""
  if (isKnownIssue == false) {
    desc = "";
  }

  $.ajax({
    url: "/testupdate",
    method: "POST",
    contentType: "application/json; charset=utf-8",
    dataType: "json",

    data: JSON.stringify({
      TestName: testName,
      TestRun: lastTestRun,
      IsKnownIssue: isKnownIssue,
      KnownIssueDescription: desc,
    }),

    statusCode: {
      200: function (xhr) {
        // TODO useful?
        // tedNotesE = button.parentNode.getElementsByClassName("tednotes");
        if (isKnownIssue == false) {
          // Clear the Known Issue input
          $("input#history-table-input-known-issue-" + testNameDown).val("");
          // TODO
          // Remove the 'known_issue' class from the result field
          // document.getElementById(resultFieldID).classList.remove('.test-known_issue');
        } else {
          // TODO
          // Add the 'known_issue' class to the result field
          // document.getElementById(resultFieldID).classList.add('.test-known_issue');
        }
      },
    },
  });
}

function reloadPage() {
  window.location.reload();
}
