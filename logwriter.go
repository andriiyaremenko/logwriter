package logw

import (
	"context"
	"io"
	"strings"
	"time"
)

// LogWriter configuration options
type LogWriterOption func() (level int, f Formatter, dateFormat string)

var (
	// LogWriter configuration constructor
	Option = func(level int, f Formatter, dateFormat string) LogWriterOption {
		return func() (int, Formatter, string) { return level, f, dateFormat }
	}

	// Default JSON LogWriter configuration
	JSONOption LogWriterOption = Option(LevelInfo, JSONFormatter, time.RFC3339)
	// Default Text LogWriter configuration
	TextOption LogWriterOption = Option(LevelInfo, TextFormatter, time.RFC3339)
)

// JSON LogWriter with default options
func JSONLogWriter(ctx context.Context, w io.Writer) io.Writer {
	return LogWriter(ctx, w, JSONOption)
}

// Text LogWriter with default options
func TextLogWriter(ctx context.Context, w io.Writer) io.Writer {
	return LogWriter(ctx, w, TextOption)
}

// Generic LogWriter constructor
func LogWriter(ctx context.Context, w io.Writer, conf LogWriterOption) io.Writer {
	level, formatter, dateTemplate := conf()
	return &logWriter{
		ctx:          ctx,
		w:            w,
		level:        level,
		dateTemplate: dateTemplate,
		formatter:    formatter,
	}
}

type logWriter struct {
	ctx context.Context
	w   io.Writer

	level        int
	dateTemplate string
	formatter    Formatter
}

func (lw *logWriter) Write(p []byte) (int, error) {
	now := time.Now().Round(time.Millisecond)
	level, message, tags := parseLog(p)

	if level < lw.level {
		return 0, nil
	}

	tags = append(getTags(lw.ctx, level), tags...)

	log := Log{
		LevelCode: level,
		Level:     getLevelText(level),
		Message:   strings.TrimRight(string(message), "\n"),
		Tags:      tags,
		Date:      now,
	}

	return lw.w.Write(lw.formatter(&log, lw.dateTemplate))
}
