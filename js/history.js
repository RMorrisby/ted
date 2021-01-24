// ############# Websocket stuff
// TODO

// ############# END Websocket stuff

// TODO other stuff?

function appendKnownIssueButtonsToLastResults() {
  var table = document.getElementById("history-table-body");

  var last = $(table).find(".child").find("td").last();

  // TODO
}

function scrollRightMax() {
  var scrollWidth = $("#history-table-body").scrollWidth;
  console.log(scrollWidth);
  $("#history-table-body").scrollLeft(scrollWidth);
}

// On page load, adorn the table with whstever extra elements we need
// JS requires this function-wrapping
window.onload = function () {
  //   appendKnownIssueButtonsToLastResults();
  scrollRightMax();
};
