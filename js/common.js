// JS file for common functions

// Returns the string with all non-alphanumeric replaced with underscores, and with all letters in lowercase
function downcaseAndUnderscore(s) {
  return s.replace(/\W+/g, "_").toLowerCase();
}
