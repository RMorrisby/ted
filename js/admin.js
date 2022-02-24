// When parts of the page are updated (e.g. after deleting a suite), it is likely that the result-count
// info (and other parts) will need to be refreshed
function resultInfoRefresh() {
  getResultCount();
  getAllTestRuns();
}

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
    if (data == 0) {
      document.getElementById("test-list").innerHTML = "";
    }
  });

  $.get("/admin/tests", function (data) {
    var json = JSON.parse(data);
    if (json == null) {
      return; // json is null if there are no tests
    }
    console.log(`Received ${json.length} tests`);

    var e = document.getElementById("test-list");
    e.innerHTML = ``;
    for (var i = 0; i < json.length; i++) {
      var obj = json[i];
      console.log("Received test " + obj.Name);
      e.innerHTML += `<li id=test-list-${downcaseAndUnderscore(obj.Name)}>${obj.Name} :: <button onclick="deleteTest('${
        obj.Name
      }')">Delete ${obj.Name}</button></li>`;
    }
  });
}

// Get the total number of suites in the store
function getSuiteCount() {
  $.get("/admin/getsuitecount", function (data) {
    document.getElementById("suitecount").textContent = " " + data;
    if (data == 0) {
      document.getElementById("suite-list").innerHTML = "";
    }
  });
  $.get("/admin/suites", function (data) {
    var json = JSON.parse(data);
    if (json == null) {
      return; // json is null if there are no suites
    }
    console.log(`Received ${json.length} suites`);

    var e = document.getElementById("suite-list");
    e.innerHTML = ``;
    for (var i = 0; i < json.length; i++) {
      var obj = json[i];
      console.log("Received suite " + obj.Name);
      e.innerHTML += `<li id=suite-list-${downcaseAndUnderscore(obj.Name)}>${
        obj.Name
      } :: <button onclick="deleteSuite('${obj.Name}')">Delete ${obj.Name}</button></li>`;
    }
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
}

// Deletes a specific test run from the store
function deleteTestRun(name) {
  $.ajax({
    url: `/testrun?testrun=${name}`,
    method: "DELETE",
    contentType: "application/json",
    success: function (data) {
      console.log(`Deleted test run ${name}`);
      e = document.getElementById(`test-run-${downcaseAndUnderscore(name)}`);
      e.parentNode.removeChild(e);

      // // Decrement the test-count
      // e = $("#testcount");
      // e.text(Number(e.text()) - 1);

      // Refresh other parts of the page
      resultInfoRefresh();
    },
    error: function (request, msg, error) {
      console.error("Failed to delete test run " + name);
      // TODO more?
    },
  });
}

// Deletes all tests from the store
function deleteAllTests() {
  $.ajax({
    url: "/admin/deletealltests",
    method: "DELETE",
    contentType: "application/json",
    success: function (data) {
      document.getElementById("testcount").textContent = " " + data;
      document.getElementById("test-list").innerHTML = "";

      // Refresh other parts of the page
      resultInfoRefresh();
    },
    error: function (request, msg, error) {
      console.error("Failed to delete all tests");
      // TODO more?
    },
  });
}

// Deletes a specific test from the store
function deleteTest(name) {
  $.ajax({
    url: `/test?test=${name}`,
    method: "DELETE",
    contentType: "application/json",
    success: function (data) {
      console.log(`Deleted test ${name}`);
      e = document.getElementById(`test-list-${downcaseAndUnderscore(name)}`);
      e.parentNode.removeChild(e);

      // Decrement the test-count
      e = $("#testcount");
      e.text(Number(e.text()) - 1);

      // Refresh other parts of the page
      resultInfoRefresh();
    },
    error: function (request, msg, error) {
      console.error("Failed to delete test " + name);
      // TODO more?
    },
  });
}

// Deletes all suites from the store
function deleteAllSuites() {
  $.ajax({
    url: "/admin/deleteallsuites",
    method: "DELETE",
    contentType: "application/json",
    success: function (data) {
      document.getElementById("suitecount").textContent = " " + data;
      document.getElementById("suite-list").innerHTML = "";

      // Refresh other parts of the page
      resultInfoRefresh();
    },
    error: function (request, msg, error) {
      console.error("Failed to delete all suites");
      // TODO more?
    },
  });
}

// Deletes a specific suite from the store
function deleteSuite(name) {
  $.ajax({
    url: `/suite?suite=${name}`,
    method: "DELETE",
    contentType: "application/json",
    success: function (data) {
      console.log(`Deleted suite ${name}`);
      e = document.getElementById(`suite-list-${downcaseAndUnderscore(name)}`);
      e.parentNode.removeChild(e);

      // Decrement the suite-count
      e = $("#suitecount");
      e.text(Number(e.text()) - 1);

      // Refresh other parts of the page
      resultInfoRefresh();
    },
    error: function (request, msg, error) {
      console.error("Failed to delete suite " + name);
      // TODO more?
    },
  });
}

// Get the names & test counts of all known test runs in the store
function getAllTestRuns() {
  console.log("Requesting all test runs...");

  $.get("/admin/getalltestruncounts", function (data) {
    console.log("Received all test runs");
    var json = JSON.parse(data);

    console.log(`Received ${json.length} test runs`);
    var e = document.getElementById("test-run-list");
    e.innerHTML = ``;
    for (var i = 0; i < json.length; i++) {
      var obj = json[i];
      console.log("Received " + obj.TestRunName + " and " + obj.Count);
      e.innerHTML += `<li id=test-run-${downcaseAndUnderscore(obj.TestRunName)}>${obj.TestRunName} :: ${
        obj.Total
      } <button onclick="deleteTestRun('${obj.TestRunName}')">Delete ${obj.TestRunName}</button></li>`;
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
