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
  var e = document.getElementById("results-table-body");

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

  e.innerHTML += `
    <tr id="${r.TestName}-${r.TestRunIdentifier}">
        <td id="categories">${r.Categories}</td>
        <td id="dir">${r.Dir}</td>
        <td id="testname">${r.TestName}</td>
        <td id="testrun">${r.TestRunIdentifier}</td>
        <td class=${testStatusClass} id="status">${r.Status}</td>
        <td id="priority">${r.Priority}</td>
        <td id="start">${startDate}</td>
        <td id="end">${endDate}</td>
        <td id="ranby">${r.RanBy}</td>
        <td id="message">${r.Message}</td>
        <td id="tedstatus">${r.TedStatus}</td>
        <td id="tednotes">${r.TedNotes}</td>
    </tr>
    `;

  // Give the TED status cell two statuses - the test status and the TED status
  // The TED status takes precedence for controlling the cell's formatting, with the test status as the backup 
  e.getElementById("tedstatus").classList.add(testStatusClass);
  e.getElementById("tedstatus").classList.add(tedStatusClass);

  // Also give the TED status cell a fixed class that can be used to get the cell (within the row)
  td.classList.add("tedstatus")

  addKnownIssueFieldsToTableRow(e);
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
