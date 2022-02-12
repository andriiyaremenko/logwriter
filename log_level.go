package logw

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	_ int = iota
	// Debug level code
	LevelDebug
	// Info level code
	LevelInfo
	// Warn level code
	LevelWarn
	// Error level code
	LevelError
	// Fatal level code
	LevelFatal
)

const (
	closingSectionHeader string = "-logw-e-\n"
	levelSectionHeader   string = "-logw-l-\n"
	tagSectionHeader     string = "-logw-t-\n"
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

var tagHeaderLen int = len(tagSectionHeader)

// Sets message level
func Level(level int) LogLevel {
	return LogLevel(strings.Join(
		[]string{
			tagSectionHeader,
			closingSectionHeader,
			fmt.Sprintf(levelSectionHeader+"%d"+closingSectionHeader, level),
		}, "",
	))
}

// Message log level
// Allows in-place tags
type LogLevel string

// Adds in-place tag with string value
func (t LogLevel) WithString(tag string, value string) LogLevel {
	return t.appendTag(fmt.Sprintf("%s\t%T\t%s", tag, value, value))
}

// Adds in-place tag with int value
func (t LogLevel) WithInt(tag string, value int) LogLevel {
	return t.appendTag(fmt.Sprintf("%s\t%T\t%d", tag, value, value))
}

// Adds in-place tag with float value
func (t LogLevel) WithFloat(tag string, value float64) LogLevel {
	return t.appendTag(fmt.Sprintf("%s\t%T\t%f", tag, value, value))
}

// Adds in-place tag with bool value
func (t LogLevel) WithBool(tag string, value bool) LogLevel {
	return t.appendTag(fmt.Sprintf("%s\t%T\t%t", tag, value, value))
}

// Adds in-place trace tag with file name and row number
// Tag key: "trace"
func (t LogLevel) WithRowNumber() LogLevel {
	return t.WithString("trace", getFileAndLine(1))
}

// Appends log message
func (t LogLevel) WithMessage(template string, v ...interface{}) string {
	return strings.Join([]string{string(t), fmt.Sprintf(template, v...)}, "")
}

func (t LogLevel) appendTag(tagSection string) LogLevel {
	return LogLevel(
		strings.Join(
			[]string{string(t[:tagHeaderLen]), tagSection, "\n", string(t[tagHeaderLen:])}, "",
		),
	)
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
