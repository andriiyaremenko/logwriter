package logw

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/andriiyaremenko/logwriter/color"
)

type Formatter func(*Log) []byte

func JSONFormatter(log *Log) []byte {
	jsonLog := make(map[string]interface{})
	jsonLog["levelCode"] = log.LevelCode
	jsonLog["level"] = log.Level
	jsonLog["date"] = log.Date.UTC().Format(time.RFC3339)
	jsonLog["message"] = log.Message

	tags := make(map[string][]interface{})
	for _, tag := range log.Tags {
		tags[tag.Key] = append(tags[tag.Key], tag.Value)
	}

	for k, v := range tags {
		jsonLog[k] = v
	}

	b, err := json.Marshal(jsonLog)
	if err != nil {
		fmt.Println(color.ColorizeText(color.ANSIColorRed, fmt.Sprintf("LogWriterJSON: failed to write log: %s", err)))

		return []byte("")
	}

	return append(b, '\n')
}
