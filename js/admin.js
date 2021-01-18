// Get the total number of results in the store
function getResultCount() {
  $.get("/admin/getresultcount", function (data) {
    document.getElementById("resultcount").textContent = " " + data;
  });
}
// Get the total number of tests in the store
function getTestCount() {
  $.get("/admin/gettestcount", function (data) {
    document.getElementById("testcount").textContent = " " + data;
  });
}

// Get the total number of suites in the store
function getSuiteCount() {
  $.get("/admin/getsuitecount", function (data) {
    document.getElementById("suitecount").textContent = " " + data;
  });
}

// Deletes all results from the store
function deleteAllResults() {
  $.ajax({
    url: "/admin/deleteallresults",
    method: "DELETE",
    contentType: "application/json",
    success: function (data) {
      document.getElementById("resultcount").textContent = " " + data;
      document.getElementById("test-run-list").innerHTML = "";
    },
    error: function (request, msg, error) {
      console.error("Failed to delete all results");
      // TODO more?
    },
  });
  // $.delete("/admin/deleteallresults", function (data) {
  //   document.getElementById("resultcount").textContent = " " + data;
  // });
  // document.getElementById("test-run-list").innerHTML = "";
}

// Deletes all tests from the store
function deleteAllTests() {
  $.ajax({
    url: "/admin/deletealltests",
    method: "DELETE",
    contentType: "application/json",
    success: function (data) {
      document.getElementById("testcount").textContent = " " + data;
      // document.getElementById("test-list").innerHTML = "";// TODO add this?
    },
    error: function (request, msg, error) {
      console.error("Failed to delete all tests");
      // TODO more?
    },
  });
  // $.delete("/admin/deletealltests", function (data) {
  //   document.getElementById("testcount").textContent = " " + data;
  // });
  // document.getElementById("test-list").innerHTML = "";
}

// Deletes all suites from the store
function deleteAllSuites() {
  $.ajax({
    url: "/admin/deleteallsuites",
    method: "DELETE",
    contentType: "application/json",
    success: function (data) {
      document.getElementById("suitecount").textContent = " " + data;
      // document.getElementById("suite-list").innerHTML = ""; // TODO add this?
    },
    error: function (request, msg, error) {
      console.error("Failed to delete all suites");
      // TODO more?
    },
  });
  // $.delete("/admin/deleteallsuites", function (data) {
  //   document.getElementById("suitecount").textContent = " " + data;
  // });
  // document.getElementById("suite-list").innerHTML = "";
}

// Get the names & test counts of all known test runs in the store
function getAllTestRuns() {
  console.log("Requesting all test runs...");

  $.get("/admin/getalltestruncounts", function (data) {
    console.log("Received all test runs");
    var json = JSON.parse(data);

    console.log(`Received ${json.length} test runs`);
    var e = document.getElementById("test-run-list");
    for (var i = 0; i < json.length; i++) {
      e.innerHTML = ``;
      for (var i = 0; i < json.length; i++) {
        var obj = json[i];
        console.log("Received " + obj.TestRunName + " and " + obj.Count);
        e.innerHTML += `<li>${obj.TestRunName} :: ${obj.Count}</li>`;
      }
    }
  });
}

// On page load, get the result-count
// JS requires this function-wrapping
window.onload = function () {
  getResultCount();
  getAllTestRuns();
  getTestCount();
  getSuiteCount();
};
