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
    const r = JSON.parse(evt.data); // thi should be a ResultStruct (JSON)
    announceLatestResult(r);
    addResultToPage(r);
  };
}

function addResultToPage(r) {
  var e = document.getElementById("results-table-body");
  e.innerHTML += `
    <tr>
        <td>${r.Category}</td>
        <td>${r.Name}</td>
        <td>${r.TestRunIdentifier}</td>
        <td>${r.Status}</td>
        <td>${r.Timestamp}</td>
        <td>${r.Message}</td>
    </tr>
    `;
}

function announceLatestResult(r) {
  var e = document.getElementById("latest-test");
  e.innerHTML =
    "<span>Latest test : " +
    r.Name +
    "    " +
    r.Status +
    " on " +
    r.Timestamp +
    "</span>";
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
        console.log(
          "ws connection failed; now trying to connect via wss to wss://" +
            wsurl
        );
        tryConnectToReload("wss://" + wsurl);
        console.log("WSS connection succeeded");
      } catch (ex2) {
        console.log("wss connection failed");
      }
    }
  } else {
    console.log(
      "Your browser does not support WebSockets, cannot connect to the Reload service."
    );
  }
} catch (ex) {
  console.error("Exception during connecting to Reload:", ex);
}

// ############# Websocket stuff

function addResultToPage(r) {
  var e = document.getElementById("results-table-body");
  e.innerHTML += `
    <tr>
        <td>${r.Category}</td>
        <td>${r.Name}</td>
        <td>${r.TestRunIdentifier}</td>
        <td>${r.Status}</td>
        <td>${r.Timestamp}</td>
        <td>${r.Message}</td>
    </tr>
    `;
}

// Get all existing results from the DB
function getAllResults() {
  console.log("Requesting all results...");

  $.get("/results", function (data) {
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
