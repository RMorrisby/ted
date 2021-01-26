package enums

type ResultStatus string

const (
	Passed        ResultStatus = "PASSED"
	Failed                     = "FAILED"
	NotRun                     = "NOT RUN"
	Unknown                    = "UNKNOWN"
	KnownIssue                 = "KNOWN ISSUE"
	Intermittent               = "INTERMITTENT"
	PassedOnRerun              = "PASSED ON RERUN"
)
