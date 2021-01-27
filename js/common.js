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



// To the table row e, append the Known Issue fields
// Takes : e : element (the tr)
//         test : Test object
//         lastTestRun : string (the value of the last test run)  
function addKnownIssueFieldsToTableRow(e, test, lastTestRun) {

  var testNameDown = downcaseAndUnderscore(test.Name)

  // Button to clear the Known Issue value
  var buttonClear = document.createElement("button");
  buttonClear.className = "known-issue-clear";
  buttonClear.id = "history-table-button-known-issue-clear-" + testNameDown;
  buttonClear.appendChild(document.createTextNode("N"));
  buttonClear.setAttribute("test", test.Name);
  buttonClear.setAttribute("testrun", lastTestRun);
  buttonClear.setAttribute("is-known-issue", false);
  $(buttonClear).on("click", function () {
    sendKnownIssueForTest(this);
  });
  // buttonClear.onclick = function () {
  //   $.ajax({
  //     url: "/testupdate",
  //     method: "POST",
  //     contentType: "application/json; charset=utf-8",
  //     dataType: "json",

  //     data: JSON.stringify({ TestName: test.Name, TestRun: lastTestRun, IsKnownIssue: false }),

  //     success: function (data) {
  //       console.log(`Updated known-issue status for test ${test.Name}`);
  //     },
  //     error: function (request, msg, error) {
  //       console.error("Failed to update test's Known Issue fields");
  //       // TODO more?
  //     },
  //   });
  // };

  // Button to set the Known Issue value
  var buttonSet = document.createElement("button");
  buttonSet.className = "known-issue-set";
  buttonSet.id = "history-table-button-known-issue-set-" + testNameDown;
  buttonSet.appendChild(document.createTextNode("Y"));
  buttonSet.setAttribute("test", test.Name);
  buttonSet.setAttribute("testrun", lastTestRun);
  buttonSet.setAttribute("is-known-issue", true);
  $(buttonSet).on("click", function () {
    sendKnownIssueForTest(this);
  });

  // Input field for the Known Issue value
  var input = document.createElement("input");
  input.className = "known-issue-input";
  input.id = "history-table-input-known-issue-" + testNameDown;
  input.value = test.KnownIssueDescription;
  input.setAttribute("test", test.Name);

  var td = document.createElement("td");
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
        // var resultFieldID = "" // TODO
        button.parent().getElementById("")
        if (isKnownIssue == false) {
          // Clear the Known Issue input
          $("input#history-table-input-known-issue-" + testNameDown).val("");
          // Add the 'known_issue' class to the result field
          // document.getElementById(resultFieldID).classList.remove('.test-known_issue');
        } else {
          // Add the 'known_issue' class to the result field
          // document.getElementById(resultFieldID).classList.add('.test-known_issue');

        }
      },
    },
  });
}
