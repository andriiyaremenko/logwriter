package logw

import "time"

// Log message
type Log struct {
	// Log level code
	LevelCode int
	// Log level
	Level string
	// Logged message
	Message string
	// Slice of tags collected from Context and in-place
	Tags []Tag
	// Log time-stamp
	Date time.Time
}
