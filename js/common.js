// JS file for common functions

// Returns the string with all whitespace replaced with underscores, and with all letters in lowercase
// TODO
function downcaseAndUnderscore(s) {
  // return s.replace(/\W+/g, "_").toLowerCase();
  // return s.replace(/\s+/g, "_").toLowerCase();
  return s.replace(/[^a-zA-Z0-9.-_=+]+/g, "_").toLowerCase();
}
