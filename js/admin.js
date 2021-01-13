// Get the total number of results in the store
function getResultCount() {
  var e = document.getElementById("resultcount");
  var xhr = new XMLHttpRequest();
  xhr.onreadystatechange = function () {
    if (xhr.readyState == 4 && xhr.status == 200) {
      e.textContent = " " + xhr.responseText;
    }
  };
  xhr.open("GET", "/admin/getcount", true);
  try {
    xhr.send();
  } catch (err) {
    /* handle error */
  }
}

// Deletes all results from the store
function deleteAllResults() {
  var e = document.getElementById("resultcount");
  var xhr = new XMLHttpRequest();
  xhr.onreadystatechange = function () {
    if (xhr.readyState == 4 && xhr.status == 200) {
      e.textContent = " " + xhr.responseText;
    }
  };
  xhr.open("POST", "/admin/deleteall", true);
  try {
    xhr.send();
  } catch (err) {
    /* handle error */
  }
}

// Get the names & test counts of all known test runs in the store
function getAllTestRuns() {
  var e = document.getElementById("test-run-list");
  var xhr = new XMLHttpRequest();
  xhr.onreadystatechange = function () {
    if (xhr.readyState == 4 && xhr.status == 200) {
      console.log("Received " + xhr.responseText);
      var json = JSON.parse(xhr.responseText);

      e.innerHTML = ``;
      for (var i = 0; i < json.length; i++) {
        var obj = json[i];
        console.log("Received " + obj.TestRunName + " and " + obj.Count);
        e.innerHTML += `<li>${obj.TestRunName} :: ${obj.Count}</li>`;
      }
    }
  };
  xhr.open("GET", "/admin/getalltestruncounts", true);
  try {
    xhr.send();
  } catch (err) {
    /* handle error */
  }
}

// On page load, get the result-count
// JS requires this function-wrapping
window.onload = function () {
  getResultCount();
  getAllTestRuns();
};
