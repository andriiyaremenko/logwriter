package logw

import (
	"bytes"
	"regexp"
	"strconv"
)

var (
	regexLevel = regexp.MustCompile(levelSection + `(?P<level>\d+)` + closingSection)
	regexTags  = regexp.MustCompile(
		tagSection + `(\w+) (string|int|float64|bool) (true|false|\d+|.+)` + closingSection,
	)
)

func parseLog(m []byte) (int, []byte, []Tag) {
	tags := []Tag{}
	if !regexLevel.Match(m) && !regexTags.Match(m) {
		return LevelInfo, m, tags
	}

	levelMatch := regexLevel.FindSubmatch(m)
	levelIndex := regexLevel.SubexpIndex("level")

	message := regexLevel.ReplaceAll(m, []byte{})

	level, err := strconv.Atoi(string(levelMatch[levelIndex]))
	if err != nil {
		level = LevelInfo
	}

	if !regexTags.Match(m) {
		message = bytes.TrimLeft(message, " ")

		return level, message, tags
	}

	tagsMatch := regexTags.FindAllSubmatch(message, -1)
	message = regexTags.ReplaceAll(message, []byte{})
	message = bytes.TrimLeft(message, " ")

	for i := range tagsMatch {
		var value interface{}
		switch string(tagsMatch[i][2]) {
		case "string":
			value = string(tagsMatch[i][3])
		case "bool":
			value, err = strconv.ParseBool(string(tagsMatch[i][3]))
			if err != nil {
				value = tagsMatch[i][3]
			}
		case "int":
			value, err = strconv.Atoi(string(tagsMatch[i][3]))
			if err != nil {
				value = tagsMatch[i][3]
			}
		case "float64":
			value, err = strconv.ParseFloat(string(tagsMatch[i][3]), 64)
			if err != nil {
				value = tagsMatch[i][3]
			}
		default:
			value = tagsMatch[i][3]
		}

		tags = append(tags, Tag{
			Key:   string(tagsMatch[i][1]),
			Value: value,
			Level: level,
		})
	}

	return level, message, tags
}

func getLevelText(level int) string {
	switch level {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	}

	if level < LevelDebug {
		return "trace"
	}

	return "fatal"
}
