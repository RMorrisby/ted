function addResultToPage(r) {
  var tbody = document.getElementById("results-table-body");

  var testStatusClass = "test-" + downcaseAndUnderscore(r.Status);
  var tedStatusClass = "test-" + downcaseAndUnderscore(r.TedStatus);

  console.log(r); // TODO remove
  // If it is absent, set the message to an empty string
  if (r.Message == null) {
    r.Message = "";
  }
  // If it is absent, set the TedNotes to an empty string
  if (r.TedNotes == null) {
    r.TedNotes = "";
  }

  // type Result struct {
  //   TestName          string
  //   SuiteName         string
  //   TestRunIdentifier string
  //   Status            string
  //   StartTimestamp    string
  //   EndTimestamp      string
  //   RanBy             string
  //   Message           string
  //   TedStatus         string
  //   TedNotes          string
  // }

  tbody.innerHTML += `
    <tr id="${r.TestName}-${r.TestRunIdentifier}">
        <td class="testname">${r.TestName}</td>
        <td class="suitename">${r.SuiteName}</td>
        <td class="testrun">${r.TestRunIdentifier}</td>
        <td class=${testStatusClass}>${r.Status}</td>
        <td class="start">${r.StartTimestamp}</td>
        <td class="end">${r.EndTimestamp}</td>
        <td class="ranby">${r.RanBy}</td>
        <td class="message">${r.Message}</td>
        <td class="tedstatus">${r.TedStatus}</td>
        <td class="tednotes">${r.TedNotes}</td>
    </tr>
    `;

  // var tr = document.createElement("tr");
  // tr.id = `${r.TestName}-${r.TestRunIdentifier}`;

  // var td = document.createElement("td");
  // td.className = "categories";
  // td.appendChild(document.createTextNode(r.Categories));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "dir";
  // td.appendChild(document.createTextNode(r.Dir));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "testname";
  // td.appendChild(document.createTextNode(r.TestName));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "testrun";
  // td.appendChild(document.createTextNode(r.TestRunIdentifier));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.classList.add(testStatusClass);
  // td.classList.add("status");
  // td.appendChild(document.createTextNode(r.Status));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "priority";
  // td.appendChild(document.createTextNode(r.Priority));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "start";
  // td.appendChild(document.createTextNode(startDate));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "end";
  // td.appendChild(document.createTextNode(endDate));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "ranby";
  // td.appendChild(document.createTextNode(r.RanBy));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "message";
  // td.appendChild(document.createTextNode(r.Message));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "tedstatus";
  // td.appendChild(document.createTextNode(r.TedStatus));
  // tr.appendChild(td);

  // // Add to the TED status cell two statuses - the test status and the TED status
  // // The TED status takes precedence for controlling the cell's formatting, with the test status as the backup
  // // td = tbody.getElementsByClassName("tedstatus")[-1];
  // td.classList.add(testStatusClass);
  // td.classList.add(tedStatusClass);
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "tednotes";
  // td.appendChild(document.createTextNode(r.TedNotes));
  // tr.appendChild(td);

  // // tr = tbody.getElementsByTagName("tr")[-1];
  // addKnownIssueFieldsToTableRow(tr, r.TestName, r.TestRunIdentifier, r.TedNotes);
  // tbody.appendChild(tr);
}

// Get all existing results from the DB
function getEntireResultTable() {
  console.log("Requesting all results...");

  $.get("/db_results_get_all", function (data) {
    console.log("Received entire table contents");
    var json = JSON.parse(data);

    console.log(`Received ${json.length} results`);
    for (var i = 0; i < json.length; i++) {
      var r = json[i];
      addResultToPage(r);
    }
  });
}

// On page load, get the existing results
// JS requires this function-wrapping
window.onload = function () {
  getEntireResultTable();
};
