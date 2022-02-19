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
	message []byte,
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
	message []byte,
) []byte {
	var sb strings.Builder

	sb.WriteByte('{')

	sb.WriteString("\"levelCode\":")
	sb.WriteString(strconv.Itoa(levelCode))
	sb.WriteByte(',')

	sb.WriteString("\"level\":")
	sb.WriteByte('"')
	sb.WriteString(level)
	sb.WriteByte('"')
	sb.WriteByte(',')

	if dateLayout != NoDate {
		sb.WriteString("\"date\":")
		sb.WriteByte('"')
		sb.WriteString(timeStamp.UTC().Format(dateLayout))
		sb.WriteByte('"')
		sb.WriteByte(',')
	}

	sb.WriteString("\"message\":")
	sb.WriteByte('"')

	sb.Write(message)
	sb.WriteByte('"')

	tagsMap := make(map[string][][]byte)
	for _, tag := range tags {
		tagsMap[tag.Key] = append(tagsMap[tag.Key], tag.Value)
	}

	for k, v := range tagsMap {
		sb.WriteByte(',')
		sb.WriteByte('"')
		sb.WriteString(k)
		sb.WriteByte('"')
		sb.WriteByte(':')
		sb.WriteByte('[')
		sb.Write(bytes.Join(v, []byte{','}))
		sb.WriteByte(']')
	}

	sb.WriteByte('}')
	sb.WriteByte('\n')

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
	message []byte,
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
	sb.WriteByte('\t')

	if dateLayout != NoDate {
		sb.WriteString(color.ColorizeText(color.ANSIColorGray, timeStamp.Format(dateLayout)))
		sb.WriteByte('\t')
	}

	tagsMap := make(map[string][][]byte)
	for _, tag := range tags {
		tagsMap[tag.Key] = append(tagsMap[tag.Key], tag.Value)
	}

	for k, values := range tagsMap {
		sb.WriteString(k)
		sb.WriteByte(':')
		sb.WriteByte('[')
		sb.Write(bytes.Join(values, []byte{','}))
		sb.WriteByte(']')

		sb.WriteByte('\t')
	}

	sb.Write(message)
	sb.WriteByte('\n')

	result := sb.String()
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 2, 2, ' ', 0)

	if _, err := w.Write([]byte(result)); err != nil {
		return []byte(result)
	}

	if err := w.Flush(); err != nil {
		return []byte(result)
	}

	return buf.Bytes()
}
