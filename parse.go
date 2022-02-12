package logw

import (
	"strconv"
	"strings"
)

func parseLog(m []byte) (int, []byte, []Tag) {
	tags := []Tag{}
	level := LevelInfo

	sections := strings.SplitN(string(m), logwHeader, 3)

	if len(sections) != 3 {
		return level, m, tags
	}

	message := sections[2]
	message = strings.TrimLeft(message, " ")

	for _, row := range strings.Split(sections[1], "\n") {
		tagSection := strings.SplitN(row, "\t", 3)
		if len(tagSection) != 3 {
			continue
		}

		var err error
		var value interface{}

		switch tagSection[1] {
		case "string":
			value = string(tagSection[2])
		case "bool":
			value, err = strconv.ParseBool(string(tagSection[2]))
			if err != nil {
				value = tagSection[2]
			}
		case "int":
			value, err = strconv.Atoi(string(tagSection[2]))
			if err != nil {
				value = tagSection[2]
			}
		case "float64":
			value, err = strconv.ParseFloat(string(tagSection[2]), 64)
			if err != nil {
				value = tagSection[2]
			}
		case "level":
			level, err = strconv.Atoi(string(tagSection[2]))
			if err != nil {
				level = LevelInfo
			}
			continue
		default:
			value = tagSection[2]
		}

		tags = append(tags, Tag{
			Key:   tagSection[0],
			Value: value,
		})
	}

	for i, tag := range tags {
		tag.Level = level
		tags[i] = tag
	}

	return level, []byte(message), tags
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
