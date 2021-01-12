
// Get the total number of results in the store
function getResultCount() {
    var e = document.getElementById("resultcount");
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4 && xhr.status == 200) {
            e.textContent = " " + xhr.responseText;
        }
    }
    xhr.open("GET", "/admin/getcount", true);
    try { xhr.send(); } catch (err) { /* handle error */ }
}

// Deletes all results from the store
function deleteAllResults() {
    var e = document.getElementById("resultcount");
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4 && xhr.status == 200) {
            e.textContent = " " + xhr.responseText;
        }
    }
    xhr.open("POST", "/admin/deleteall", true);
    try { xhr.send(); } catch (err) { /* handle error */ }
}