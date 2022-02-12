package logw

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

var (
	regexLevel = regexp.MustCompile(levelSectionHeader + `(?P<level>\d+)` + closingSectionHeader)
	regexTags  = regexp.MustCompile(
		tagSectionHeader + `(?P<tags>(.|\n)*)` + closingSectionHeader,
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
	message = regexTags.ReplaceAll(message, []byte{})
	message = bytes.TrimLeft(message, " ")

	level, err := strconv.Atoi(string(levelMatch[levelIndex]))
	if err != nil {
		level = LevelInfo
	}

	tagsMatch := regexTags.FindSubmatch(m)
	tagsIndex := regexTags.SubexpIndex("tags")

	if tagsIndex < 0 {
		return level, message, tags
	}

	for _, row := range strings.Split(string(tagsMatch[tagsIndex]), "\n") {
		tagSection := strings.SplitN(row, "\t", 3)
		if len(tagSection) != 3 {
			continue
		}

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
		default:
			value = tagSection[2]
		}

		tags = append(tags, Tag{
			Key:   tagSection[0],
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
