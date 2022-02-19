package logw

import (
	"bytes"
	"context"
	"io"
	"time"
)

// LogWriter configuration options
type LogWriterOption func() (level int, f Formatter, dateFormat string)

var (
	// LogWriter configuration constructor
	Option = func(level int, f Formatter, dateFormat string) LogWriterOption {
		return func() (int, Formatter, string) { return level, f, dateFormat }
	}
	// Without time-stamp option
	NoTimeStampOption = func(level int, f Formatter) LogWriterOption {
		return Option(level, f, NoDate)
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
	loggingLevel, formatter, dateTemplate := conf()

	return &logWriter{
		write: func(p []byte) (int, error) {
			now := time.Now().Round(time.Millisecond)
			level, message, tags := parseLog(p)

			if level < loggingLevel {
				return 0, nil
			}

			tags = append(getTags(ctx, level), tags...)

			return w.Write(
				formatter(
					getLevelText(level),
					level,
					tags,
					now,
					dateTemplate,
					bytes.TrimRight(message, "\n"),
				),
			)
		},
	}
}

type logWriter struct {
	write func(p []byte) (int, error)
}

func (w *logWriter) Write(p []byte) (int, error) {
	return w.write(p)
}
