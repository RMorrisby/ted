// ############# Websocket stuff

/**
 * Tries to connect to the reload service and start listening to reload events.
 *
 * @function tryConnectToReload
 * @public
 */
function tryConnectToReload(address) {
  var conn = new WebSocket(address);

  conn.onclose = function () {
    setTimeout(function () {
      tryConnectToReload(address);
    }, 2000);
  };

  conn.onmessage = function (evt) {
    const r = JSON.parse(evt.data); // this should be a ResultStruct (JSON)
    announceLatestResult(r);
    addResultToPage(r);
  };
}

function announceLatestResult(r) {
  var e = document.getElementById("latest-test");
  e.innerHTML = "<span>Latest test : " + r.TestName + "    " + r.Status + " on " + r.EndTimestamp + "</span>";
}

try {
  if (window.WebSocket === undefined) {
    $("#container").append("Your browser does not support WebSockets");
    //return;
  } else if (window["WebSocket"]) {
    // The reload endpoint is hosted on a statically defined port.
    try {
      if (window.location.hostname == "localhost") {
        var wsurl = "localhost:8080/datareload";
      } else {
        var wsurl = window.location.hostname + "/datareload";
      }
      console.log("Trying to connect websocket to ws://" + wsurl);
      tryConnectToReload("ws://" + wsurl);
      console.log("WS connection succeeded");
    } catch (ex) {
      // If an exception is thrown, that means that we couldn't connect to to WebSockets because of mixed content
      // security restrictions, so we try to connect using wss.
      try {
        console.log("ws connection failed; now trying to connect via wss to wss://" + wsurl);
        tryConnectToReload("wss://" + wsurl);
        console.log("WSS connection succeeded");
      } catch (ex2) {
        console.log("wss connection failed");
      }
    }
  } else {
    console.log("Your browser does not support WebSockets, cannot connect to the Reload service.");
  }
} catch (ex) {
  console.error("Exception during connecting to Reload:", ex);
}

// ############# END Websocket stuff

function addResultToPage(r) {
  var tbody = document.getElementById("results-table-body");

  if (r.StartTimestamp != null) {
    // toISOString should yield a date in this format : 2021-01-17T19:41:00.000Z
    // We want 2021-01-17 19:41
    // TODO This needs to handle non-GMT timestamps properly - we're not displaying the timezone, so to
    // the user it looks like a local time
    // The page also needs to warn / declare this
    // Incredibly, JS doesn't have any handling for format-strings. So we have to brute-force this somewhat.
    var startDate = new Date(r.StartTimestamp).toISOString().replace(/(T|Z)/g, " ").slice(0, 16);
  } else {
    var startDate = null;
  }

  if (r.EndTimestamp != null) {
    var endDate = new Date(r.EndTimestamp).toISOString().replace(/(T|Z)/g, " ").slice(0, 16);
  } else {
    var endDate = null;
  }

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
  // tbody.innerHTML += `
  //   <tr id="${r.TestName}-${r.TestRunIdentifier}">
  //       <td class="categories">${r.Categories}</td>
  //       <td class="dir">${r.Dir}</td>
  //       <td class="testname">${r.TestName}</td>
  //       <td class="testrun">${r.TestRunIdentifier}</td>
  //       <td class=${testStatusClass}>${r.Status}</td>
  //       <td class="priority">${r.Priority}</td>
  //       <td class="start">${startDate}</td>
  //       <td class="end">${endDate}</td>
  //       <td class="ranby">${r.RanBy}</td>
  //       <td class="message">${r.Message}</td>
  //       <td class="tedstatus">${r.TedStatus}</td>
  //       <td class="tednotes">${r.TedNotes}</td>
  //   </tr>
  //   `;

  var tr = document.createElement("tr");
  tr.id = `${r.TestName}-${r.TestRunIdentifier}`;

  var td = document.createElement("td");
  td.className = "categories";
  td.appendChild(document.createTextNode(r.Categories));
  tr.appendChild(td);

  var td = document.createElement("td");
  td.className = "dir";
  td.appendChild(document.createTextNode(r.Dir));
  tr.appendChild(td);

  var td = document.createElement("td");
  td.className = "testname";
  td.appendChild(document.createTextNode(r.TestName));
  tr.appendChild(td);

  var td = document.createElement("td");
  td.className = "testrun";
  td.appendChild(document.createTextNode(r.TestRunIdentifier));
  tr.appendChild(td);

  var td = document.createElement("td");
  td.classList.add(testStatusClass);
  td.classList.add("status");
  td.appendChild(document.createTextNode(r.Status));
  tr.appendChild(td);

  var td = document.createElement("td");
  td.className = "priority";
  td.appendChild(document.createTextNode(r.Priority));
  tr.appendChild(td);

  var td = document.createElement("td");
  td.className = "start";
  td.appendChild(document.createTextNode(startDate));
  tr.appendChild(td);

  var td = document.createElement("td");
  td.className = "end";
  td.appendChild(document.createTextNode(endDate));
  tr.appendChild(td);

  var td = document.createElement("td");
  td.className = "ranby";
  td.appendChild(document.createTextNode(r.RanBy));
  tr.appendChild(td);

  var td = document.createElement("td");
  td.className = "message";
  td.appendChild(document.createTextNode(r.Message));
  tr.appendChild(td);

  var td = document.createElement("td");
  td.className = "tedstatus";
  var text = r.TedStatus;
  if (r.TedNotes != null && r.TedNotes != "") {
    text = r.TedNotes;
  }
  td.appendChild(document.createTextNode(text));
  tr.appendChild(td);

  // Add to the TED status cell two statuses - the test status and the TED status
  // The TED status takes precedence for controlling the cell's formatting, with the test status as the backup
  td.classList.add(testStatusClass);
  td.classList.add(tedStatusClass);
  tr.appendChild(td);

  var td = document.createElement("td");
  td.className = "tednotes";
  td.appendChild(document.createTextNode(r.TedNotes));
  tr.appendChild(td);

  // tr = tbody.getElementsByTagName("tr")[-1];
  // var knownIssueText
  addKnownIssueFieldsToTableRow(tr, r.TestName, r.TestRunIdentifier, r.TedNotes);
  tbody.appendChild(tr);
}

// Get all existing results from the DB
function getAllResults() {
  console.log("Requesting all results...");

  var x = $(location).attr("search");

  $.get("/results" + x, function (data) {
    console.log("Received all results");
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
  getAllResults();
};
