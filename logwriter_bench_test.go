package logw_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	logw "github.com/andriiyaremenko/logwriter"
)

func BenchmarkLogLevelCompositionInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = logw.Info.WithInt("attempt", i)
	}
}

func BenchmarkLogLevelCompositionBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = logw.Info.WithBool("attempting", true)
	}
}

func BenchmarkLogLevelCompositionFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = logw.Info.WithFloat("someFloat", 3.4)
	}
}

func BenchmarkLogLevelCompositionString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = logw.Info.WithString("greeting", "Hello World")
	}
}

func BenchmarkLogLevelCompositionTrace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = logw.Info.WithTrace()
	}
}

func BenchmarkLogLevelCompositionError(b *testing.B) {
	err := errors.New("some error")
	for i := 0; i < b.N; i++ {
		_ = logw.Info.Error(err)
	}
}

func BenchmarkLogLevelCompositionAllTags(b *testing.B) {
	err := errors.New("some error")
	for i := 0; i < b.N; i++ {
		_ = logw.Info.
			WithInt("attempt", i).
			WithBool("attempting", true).
			WithFloat("someFloat", 3.4).
			Error(err).
			WithString("greeting", "Hello World")
	}
}

func BenchmarkLogLevelCompositionWithFormattedMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = logw.Info.WithMessage("this going to be fun: %d", i)
	}
}

func BenchmarkLogLevelCompositionWithAllTagsAndFormattedMessage(b *testing.B) {
	err := errors.New("some error")
	for i := 0; i < b.N; i++ {
		_ = logw.Info.
			WithInt("attempt", i).
			WithBool("attempting", true).
			WithFloat("someFloat", 3.4).
			WithString("greeting", "Hello World").
			Error(err).
			WithMessage("this going to be fun: %d", i)
	}
}

func BenchmarkLogFormatterWrite(b *testing.B) {
	ctx := context.TODO()
	f := func(
		level string,
		levelCode int,
		tags []logw.Tag,
		timeStamp time.Time,
		dateLayout string,
		message []byte,
	) []byte {
		return []byte{}
	}
	writer := logw.LogWriter(ctx, io.Discard, logw.Option(logw.LevelInfo, f, logw.NoDate))
	m := []byte(
		logw.Info.
			WithInt("attempt", 1).
			WithBool("attempting", true).
			WithFloat("someFloat", 3.4).
			WithString("greeting", "Hello World").
			Error(errors.New("some error")).
			WithMessage("this going to be fun: %d", 1),
	)

	for i := 0; i < b.N; i++ {
		_, _ = writer.Write(m)
	}
}

func BenchmarkLogWriterJSONOnlyMessage(b *testing.B) {
	ctx := context.TODO()
	writer := logw.JSONLogWriter(ctx, io.Discard)

	for i := 0; i < b.N; i++ {
		_, _ = writer.Write([]byte("this going to be fun "))
	}
}

func BenchmarkLogWriterTextOnlyMessage(b *testing.B) {
	ctx := context.TODO()
	writer := logw.TextLogWriter(ctx, io.Discard)

	for i := 0; i < b.N; i++ {
		_, _ = writer.Write([]byte("this going to be fun"))
	}
}

func BenchmarkLogWriterJSONLevelWithFormatedMessage(b *testing.B) {
	ctx := context.TODO()
	writer := logw.JSONLogWriter(ctx, io.Discard)

	m := logw.Info.WithMessage("this going to be fun: %d", 1)
	for i := 0; i < b.N; i++ {
		_, _ = writer.Write([]byte(m))
	}
}

func BenchmarkLogWriterTextLevelWithFormattedMessage(b *testing.B) {
	ctx := context.TODO()
	writer := logw.TextLogWriter(ctx, io.Discard)

	m := logw.Info.WithMessage("this going to be fun: %d", 1)
	for i := 0; i < b.N; i++ {
		_, _ = writer.Write([]byte(m))
	}
}

func BenchmarkLogWriterJSONContextTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")
	writer := logw.JSONLogWriter(ctx, io.Discard)

	for i := 0; i < b.N; i++ {
		_, _ = writer.Write([]byte("this going to be fun "))
	}
}

func BenchmarkLogWriterTextContextTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")
	writer := logw.TextLogWriter(ctx, io.Discard)

	for i := 0; i < b.N; i++ {
		_, _ = writer.Write([]byte("this going to be fun "))
	}
}

func BenchmarkLogWriterJSONAllInPlaceTags(b *testing.B) {
	ctx := context.TODO()
	writer := logw.JSONLogWriter(ctx, io.Discard)
	m := logw.Info.
		WithInt("attempt", 1).
		WithBool("attempting", true).
		WithFloat("someFloat", 3.4).
		WithString("greeting", "Hello World").
		Error(errors.New("some error")).
		WithMessage("this going to be fun: %d", 1)
	for i := 0; i < b.N; i++ {
		_, _ = writer.Write([]byte(m))
	}
}

func BenchmarkLogWriterTextAllInPlaceTags(b *testing.B) {
	ctx := context.TODO()
	writer := logw.TextLogWriter(ctx, io.Discard)
	m :=
		logw.Info.
			WithInt("attempt", 1).
			WithBool("attempting", true).
			WithFloat("someFloat", 3.4).
			WithString("greeting", "Hello World").
			Error(errors.New("some error")).
			WithMessage("this going to be fun: %d", 1)
	for i := 0; i < b.N; i++ {
		_, _ = writer.Write([]byte(m))
	}
}

func BenchmarkLogWriterJSONAllTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")
	writer := logw.JSONLogWriter(ctx, io.Discard)

	m := logw.Info.
		WithInt("attempt", 1).
		WithBool("attempting", true).
		WithFloat("someFloat", 3.4).
		WithString("greeting", "Hello World").
		Error(errors.New("some error")).
		WithMessage("this going to be fun: %d", 1)
	for i := 0; i < b.N; i++ {
		_, _ = writer.Write([]byte(m))
	}
}

func BenchmarkLogWriterTextAllTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")
	writer := logw.TextLogWriter(ctx, io.Discard)

	m := logw.Info.
		WithInt("attempt", 1).
		WithBool("attempting", true).
		WithFloat("someFloat", 3.4).
		WithString("greeting", "Hello World").
		Error(errors.New("some error")).
		WithMessage("this going to be fun: %d", 1)
	for i := 0; i < b.N; i++ {
		_, _ = writer.Write([]byte(m))
	}
}
