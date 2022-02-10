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
	Debug Message = Level(LevelDebug)
	// Sets Info message level
	Info Message = Level(LevelInfo)
	// Sets Warn message level
	Warn Message = Level(LevelWarn)
	// Sets Error message level
	Error Message = Level(LevelError)
	// Sets Fatal message level
	Fatal Message = Level(LevelFatal)
)

// Sets message level
func Level(level int) Message {
	return Message(fmt.Sprintf(levelSection+"%d"+closingSection, level))
}

// Allows in-place tags
type Message string

// Adds in-place tag with string value
func (t Message) WithString(tag string, value string) Message {
	return Message(
		fmt.Sprintf(
			string(t)+tagSection+"%s %T %s"+closingSection,
			tag,
			value,
			value,
		),
	)
}

// Adds in-place tag with int value
func (t Message) WithInt(tag string, value int) Message {
	return Message(
		fmt.Sprintf(
			string(t)+tagSection+"%s %T %d"+closingSection,
			tag,
			value,
			value,
		),
	)
}

// Adds in-place tag with float value
func (t Message) WithFloat(tag string, value float64) Message {
	return Message(
		fmt.Sprintf(
			string(t)+tagSection+"%s %T %f"+closingSection,
			tag,
			value,
			value,
		),
	)
}

// Adds in-place tag with bool value
func (t Message) WithBool(tag string, value bool) Message {
	return Message(
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
func (t Message) WithRowNumber() Message {
	value := getFileAndLine(1)
	return Message(
		fmt.Sprintf(
			string(t)+tagSection+"%s %T %s"+closingSection,
			"trace",
			value,
			value,
		),
	)
}

// Appends log message
func (t Message) WithMessage(template string, v ...interface{}) string {
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
