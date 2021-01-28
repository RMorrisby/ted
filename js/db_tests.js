function addTestToPage(t) {
  var tbody = document.getElementById("tests-table-body");

  console.log(t); // TODO remove
  // If it is absent, set the message to an empty string
  if (t.Message == null) {
    t.Message = "";
  }
  // If it is absent, set the TedNotes to an empty string
  if (t.TedNotes == null) {
    t.TedNotes = "";
  }

  // type Test struct {
  //   Name                  string
  //   Dir                   string
  //   Priority              int
  //   Categories            string // pipe-separated string
  //   Description           string
  //   Notes                 string
  //   Owner                 string
  //   IsKnownIssue          bool
  //   KnownIssueDescription string
  // }

  tbody.innerHTML += `
    <tr id="${t.Name}-${t.Dir}">
        <td class="testname">${t.Name}</td>
        <td class="dir">${t.Dir}</td>
        <td class="priority">${t.Priority}</td>
        <td class="categories">${t.Categories}</td>
        <td class="description">${t.Description}</td>
        <td class="notes">${t.Notes}</td>
        <td class="owner">${t.Owner}</td>
        <td class="isknownissue">${t.IsKnownIssue}</td>
        <td class="knownissuedescription">${t.KnownIssueDescription}</td>
    </tr>
    `;

  // var tr = document.createElement("tr");
  // tr.id = `${t.TestName}-${t.TestRunIdentifier}`;

  // var td = document.createElement("td");
  // td.className = "categories";
  // td.appendChild(document.createTextNode(t.Categories));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "dir";
  // td.appendChild(document.createTextNode(t.Dir));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "testname";
  // td.appendChild(document.createTextNode(t.TestName));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "testrun";
  // td.appendChild(document.createTextNode(t.TestRunIdentifier));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.classList.add(testStatusClass);
  // td.classList.add("status");
  // td.appendChild(document.createTextNode(t.Status));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "priority";
  // td.appendChild(document.createTextNode(t.Priority));
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
  // td.appendChild(document.createTextNode(t.RanBy));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "message";
  // td.appendChild(document.createTextNode(t.Message));
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "tedstatus";
  // td.appendChild(document.createTextNode(t.TedStatus));
  // tr.appendChild(td);

  // // Add to the TED status cell two statuses - the test status and the TED status
  // // The TED status takes precedence for controlling the cell's formatting, with the test status as the backup
  // // td = tbody.getElementsByClassName("tedstatus")[-1];
  // td.classList.add(testStatusClass);
  // td.classList.add(tedStatusClass);
  // tr.appendChild(td);

  // var td = document.createElement("td");
  // td.className = "tednotes";
  // td.appendChild(document.createTextNode(t.TedNotes));
  // tr.appendChild(td);

  // // tr = tbody.getElementsByTagName("tr")[-1];
  // addKnownIssueFieldsToTableRow(tr, t.TestName, t.TestRunIdentifier, t.TedNotes);
  // tbody.appendChild(tr);
}

// Get all existing results from the DB
function getEntireTestTable() {
  console.log("Requesting all tests...");

  $.get("/db_tests_get_all", function (data) {
    console.log("Received entire table contents");
    var json = JSON.parse(data);

    console.log(`Received ${json.length} tests`);
    for (var i = 0; i < json.length; i++) {
      var t = json[i];
      addTestToPage(t);
    }
  });
}

// On page load, get the existing results
// JS requires this function-wrapping
window.onload = function () {
  getEntireTestTable();
};
