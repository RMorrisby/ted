package enums

// Enums for the type-values that are expected in Status DB table, and enums for some of the values

type Status string

// names
const (
	StatusNamePause = "Testrun Pause"
)

// types
const (
	Pause = "PAUSE"
)

// values
const (
	Paused   = "PAUSED"
	Unpaused = "UNPAUSED"
)
