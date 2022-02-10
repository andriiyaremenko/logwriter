package logw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/andriiyaremenko/logwriter/color"
)

// Date layout to exclude logw time-stamp from log
const NoDate = "NO_DATE"

// Log message formatter
type Formatter func(log *Log, dateLayout string) []byte

// JSON message formatter.
// Has format of:
//  { "date": string|optional, "level":string, "levelCode":int, "message":string }
func JSONFormatter(log *Log, dateLayout string) []byte {
	jsonLog := make(map[string]interface{})
	jsonLog["levelCode"] = log.LevelCode
	jsonLog["level"] = log.Level
	if dateLayout != NoDate {
		jsonLog["date"] = log.Date.UTC().Format(dateLayout)
	}
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

// Text message formatter.
// Has format of:
//  level  ?time-stamp  tag-key:tag-value  message
func TextFormatter(log *Log, dateTemplate string) []byte {
	var sb strings.Builder
	levelColor := color.GetLevelColor(log.LevelCode)
	adjust := func(s string) string {
		if s == "info" || s == "warn" {
			return " " + s
		}

		return s
	}

	sb.WriteString(color.ColorizeText(levelColor, adjust(log.Level)))
	sb.WriteString("\t")

	if dateTemplate != NoDate {
		sb.WriteString(color.ColorizeText(color.ANSIColorGray, log.Date.Format(dateTemplate)))
		sb.WriteString("\t")
	}

	for _, tag := range log.Tags {
		sb.WriteString(tag.Key)
		sb.WriteString(":")
		sb.WriteString(fmt.Sprintf("%v", tag.Value))

		sb.WriteString("\t")
	}

	sb.WriteString(strings.TrimRight(log.Message, "\n"))
	sb.WriteString("\n")

	message := sb.String()
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 2, 2, ' ', 0)

	fmt.Fprint(w, message)
	w.Flush()

	return buf.Bytes()
}
