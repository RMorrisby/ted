// Get the total number of results in the store
function getResultCount() {
  $.get("/admin/getcount", function (data) {
    document.getElementById("resultcount").textContent = " " + data;
  });
}

// Deletes all results from the store
function deleteAllResults() {
  $.post("/admin/deleteall", function (data) {
    document.getElementById("resultcount").textContent = " " + data;
  });
  document.getElementById("test-run-list").innerHTML = "";
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
};
