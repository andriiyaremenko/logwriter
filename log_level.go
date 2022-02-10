package logw

import (
	"fmt"
	"runtime"
)

const (
	closingSection = `
-------------------------------------------------------------------------------------
`
	levelSection string = `
-----------------------------------logwriter-level-----------------------------------
level `
	tagSection string = `
-----------------------------------logwriter--tags-----------------------------------
`
)

var (
	// Sets Debug message level
	Debug LogLevel = Level(LevelDebug)
	// Sets Info message level
	Info LogLevel = Level(LevelInfo)
	// Sets Warn message level
	Warn LogLevel = Level(LevelWarn)
	// Sets Error message level
	Error LogLevel = Level(LevelError)
	// Sets Fatal message level
	Fatal LogLevel = Level(LevelFatal)
)

// Sets message level
func Level(level int) LogLevel {
	return LogLevel(fmt.Sprintf(levelSection+"%d"+closingSection, level))
}

// Message log level
// Allows in-place tags
type LogLevel string

// Adds in-place tag with string value
func (t LogLevel) WithString(tag string, value string) LogLevel {
	return LogLevel(
		fmt.Sprintf(
			string(t)+tagSection+"%s %T %s"+closingSection,
			tag,
			value,
			value,
		),
	)
}

// Adds in-place tag with int value
func (t LogLevel) WithInt(tag string, value int) LogLevel {
	return LogLevel(
		fmt.Sprintf(
			string(t)+tagSection+"%s %T %d"+closingSection,
			tag,
			value,
			value,
		),
	)
}

// Adds in-place tag with float value
func (t LogLevel) WithFloat(tag string, value float64) LogLevel {
	return LogLevel(
		fmt.Sprintf(
			string(t)+tagSection+"%s %T %f"+closingSection,
			tag,
			value,
			value,
		),
	)
}

// Adds in-place tag with bool value
func (t LogLevel) WithBool(tag string, value bool) LogLevel {
	return LogLevel(
		fmt.Sprintf(
			string(t)+tagSection+"%s %T %t"+closingSection,
			tag,
			value,
			value,
		),
	)
}

// Adds in-place trace tag with file name and row number
// Tag key: "trace"
func (t LogLevel) WithRowNumber() LogLevel {
	value := getFileAndLine(1)
	return LogLevel(
		fmt.Sprintf(
			string(t)+tagSection+"%s %T %s"+closingSection,
			"trace",
			value,
			value,
		),
	)
}

// Appends log message
func (t LogLevel) WithMessage(template string, v ...interface{}) string {
	return string(t) + fmt.Sprintf(template, v...)
}

func getFileAndLine(calldepth int) string {
	_, file, line, ok := runtime.Caller(calldepth + 1)

	if !ok {
		file = "???"
		line = 0
	}

	short := file

	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	return fmt.Sprintf("[%s %d]", short, line)
}
