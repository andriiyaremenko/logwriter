package logw

import (
	"encoding/json"
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
		tagSection := strings.SplitN(row, "\t", 2)
		if len(tagSection) != 2 {
			continue
		}

		var err error

		if tagSection[0] == "_level" {
			level, err = strconv.Atoi(string(tagSection[1]))
			if err != nil {
				level = LevelInfo
			}

			continue
		}

		tags = append(tags, Tag{
			Key:   tagSection[0],
			Value: json.RawMessage(tagSection[1]),
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
