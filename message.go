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
	Debug Message = Level(levelDebug)
	Info  Message = Level(levelInfo)
	Warn  Message = Level(levelWarn)
	Error Message = Level(levelError)
	Fatal Message = Level(levelFatal)
)

func Level(level int) Message {
	return Message(fmt.Sprintf(levelSection+"%d"+closingSection, level))
}

type Message string

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

func (t Message) WithRowNumber() Message {
	value := getFileAndLine(1)
	return Message(
		fmt.Sprintf(
			string(t)+tagSection+"%s %T %s"+closingSection,
			"rownumber",
			value,
			value,
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

	return fmt.Sprintf("%s: %d", short, line)
}
