package logw

import (
	"bytes"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/andriiyaremenko/logwriter/color"
)

// Date layout to exclude logw time-stamp from log
const NoDate = "NO_DATE"

// Log message formatter
type Formatter func(log *Log, dateLayout string) []byte

// JSON message formatter
// Has format of:
//  { "date": string|optional, "level":string, "levelCode":int, "message":string }
func JSONFormatter(log *Log, dateLayout string) []byte {
	var sb strings.Builder

	sb.WriteString("{")

	sb.WriteString("\"levelCode\":")
	sb.WriteString(strconv.Itoa(log.LevelCode))
	sb.WriteString(",")

	sb.WriteString("\"level\":")
	sb.WriteString("\"")
	sb.WriteString(log.Level)
	sb.WriteString("\"")
	sb.WriteString(",")

	if dateLayout != NoDate {
		sb.WriteString("\"date\":")
		sb.WriteString("\"")
		sb.WriteString(log.Date.UTC().Format(dateLayout))
		sb.WriteString("\"")
		sb.WriteString(",")
	}

	sb.WriteString("\"message\":")
	sb.WriteString("\"")

	sb.WriteString(log.Message)
	sb.WriteString("\"")

	tags := make(map[string][][]byte)
	for _, tag := range log.Tags {
		tags[tag.Key] = append(tags[tag.Key], tag.Value)
	}

	for k, v := range tags {
		sb.WriteString(",")
		sb.WriteString("\"")
		sb.WriteString(k)
		sb.WriteString("\"")
		sb.WriteString(":")
		sb.WriteString("[")
		sb.Write(bytes.Join(v, []byte(",")))
		sb.WriteString("]")
	}

	sb.WriteString("}\n")

	return []byte(sb.String())
}

// Text message formatter
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

	tags := make(map[string][][]byte)
	for _, tag := range log.Tags {
		tags[tag.Key] = append(tags[tag.Key], tag.Value)
	}

	for k, values := range tags {
		sb.WriteString(k)
		sb.WriteString(":")
		sb.WriteString("[")
		sb.WriteString(string(bytes.Join(values, []byte(","))))
		sb.WriteString("]")

		sb.WriteString("\t")
	}

	sb.WriteString(strings.TrimRight(log.Message, "\n"))
	sb.WriteString("\n")

	message := sb.String()
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 2, 2, ' ', 0)

	if _, err := w.Write([]byte(message)); err != nil {
		return []byte(message)
	}

	if err := w.Flush(); err != nil {
		return []byte(message)
	}

	return buf.Bytes()
}
