package logw

import (
	"bytes"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/andriiyaremenko/logwriter/color"
)

// Date layout to exclude logw time-stamp from log
const NoDate = "NO_DATE"

// Log message formatter
type Formatter func(
	level string,
	levelCode int,
	tags []Tag,
	timeStamp time.Time,
	dateLayout string,
	message string,
) []byte

// JSON message formatter
// Has format of:
//  { "date": string|optional, "level":string, "levelCode":int, "message":string }
func JSONFormatter(
	level string,
	levelCode int,
	tags []Tag,
	timeStamp time.Time,
	dateLayout string,
	message string,
) []byte {
	var sb strings.Builder

	sb.WriteString("{")

	sb.WriteString("\"levelCode\":")
	sb.WriteString(strconv.Itoa(levelCode))
	sb.WriteString(",")

	sb.WriteString("\"level\":")
	sb.WriteString("\"")
	sb.WriteString(level)
	sb.WriteString("\"")
	sb.WriteString(",")

	if dateLayout != NoDate {
		sb.WriteString("\"date\":")
		sb.WriteString("\"")
		sb.WriteString(timeStamp.UTC().Format(dateLayout))
		sb.WriteString("\"")
		sb.WriteString(",")
	}

	sb.WriteString("\"message\":")
	sb.WriteString("\"")

	sb.WriteString(message)
	sb.WriteString("\"")

	tagsMap := make(map[string][][]byte)
	for _, tag := range tags {
		tagsMap[tag.Key] = append(tagsMap[tag.Key], tag.Value)
	}

	for k, v := range tagsMap {
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
func TextFormatter(
	level string,
	levelCode int,
	tags []Tag,
	timeStamp time.Time,
	dateLayout string,
	message string,
) []byte {
	var sb strings.Builder
	levelColor := color.GetLevelColor(levelCode)
	adjust := func(s string) string {
		if s == "info" || s == "warn" {
			return " " + s
		}

		return s
	}

	sb.WriteString(color.ColorizeText(levelColor, adjust(level)))
	sb.WriteString("\t")

	if dateLayout != NoDate {
		sb.WriteString(color.ColorizeText(color.ANSIColorGray, timeStamp.Format(dateLayout)))
		sb.WriteString("\t")
	}

	tagsMap := make(map[string][][]byte)
	for _, tag := range tags {
		tagsMap[tag.Key] = append(tagsMap[tag.Key], tag.Value)
	}

	for k, values := range tagsMap {
		sb.WriteString(k)
		sb.WriteString(":")
		sb.WriteString("[")
		sb.WriteString(string(bytes.Join(values, []byte(","))))
		sb.WriteString("]")

		sb.WriteString("\t")
	}

	sb.WriteString(strings.TrimRight(message, "\n"))
	sb.WriteString("\n")

	message = sb.String()
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
