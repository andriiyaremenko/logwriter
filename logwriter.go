package logw

import (
	"context"
	"io"
	"strings"
	"time"
)

func JSONLogWriter(ctx context.Context, w io.Writer) io.Writer {
	return LogWriter(ctx, w, JSONFormatter)
}

func TextLogWriter(ctx context.Context, w io.Writer) io.Writer {
	return LogWriter(ctx, w, TextFormatter)
}

func LogWriter(ctx context.Context, w io.Writer, formatter Formatter) io.Writer {
	return &logWriter{ctx: ctx, w: w, formatter: formatter}
}

type logWriter struct {
	ctx       context.Context
	w         io.Writer
	formatter Formatter
}

func (lw *logWriter) Write(p []byte) (int, error) {
	now := time.Now().Round(time.Millisecond)
	level, message, tags := parseLog(p)
	tags = append(getTags(lw.ctx, level), tags...)

	log := Log{
		LevelCode: level,
		Level:     getLevelText(level),
		Message:   strings.TrimRight(string(message), "\n"),
		Tags:      tags,
		Date:      now,
	}

	return lw.w.Write(lw.formatter(&log))
}
