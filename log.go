package logw

import "time"

type Log struct {
	LevelCode int
	Level     string
	Message   string
	Tags      []Tag
	Date      time.Time
}
