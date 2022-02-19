package logw_test

import (
	"context"
	"io"
	"testing"
	"time"

	logw "github.com/andriiyaremenko/logwriter"
)

func BenchmarkLogLevelCompositionNoMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = logw.Info.
			WithInt("attempt", i).
			WithBool("attempting", true).
			WithFloat("someFloat", 3.4).
			WithString("greeting", "Hello World")
	}
}

func BenchmarkLogLevelCompositionWithMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = logw.Info.
			WithInt("attempt", i).
			WithBool("attempting", true).
			WithFloat("someFloat", 3.4).
			WithString("greeting", "Hello World").
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
		message string,
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
			WithMessage("this going to be fun: %d", 1),
	)

	for i := 0; i < b.N; i++ {
		if _, err := writer.Write(m); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWriterJSON(b *testing.B) {
	ctx := context.TODO()
	writer := logw.JSONLogWriter(ctx, io.Discard)

	for i := 0; i < b.N; i++ {
		if _, err := writer.Write([]byte("this going to be fun ")); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWriterText(b *testing.B) {
	ctx := context.TODO()
	writer := logw.TextLogWriter(ctx, io.Discard)

	for i := 0; i < b.N; i++ {
		if _, err := writer.Write([]byte("this going to be fun")); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWriterJSONLevel(b *testing.B) {
	ctx := context.TODO()
	writer := logw.JSONLogWriter(ctx, io.Discard)

	for i := 0; i < b.N; i++ {
		_, err := writer.Write([]byte(logw.Info.WithMessage("this going to be fun: %d", i)))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWriterTextLevel(b *testing.B) {
	ctx := context.TODO()
	writer := logw.TextLogWriter(ctx, io.Discard)

	for i := 0; i < b.N; i++ {
		_, err := writer.Write([]byte(logw.Info.WithMessage("this going to be fun: %d", i)))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWriterJSONContextTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")
	writer := logw.JSONLogWriter(ctx, io.Discard)

	for i := 0; i < b.N; i++ {
		_, err := writer.Write([]byte("this going to be fun "))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWriterTextContextTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")
	writer := logw.TextLogWriter(ctx, io.Discard)

	for i := 0; i < b.N; i++ {
		_, err := writer.Write([]byte("this going to be fun "))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWriterJSONInPlaceTags(b *testing.B) {
	ctx := context.TODO()
	writer := logw.JSONLogWriter(ctx, io.Discard)

	for i := 0; i < b.N; i++ {
		_, err := writer.Write(
			[]byte(
				logw.Info.
					WithInt("attempt", i).
					WithBool("attempting", true).
					WithFloat("someFloat", 3.4).
					WithString("greeting", "Hello World").
					WithMessage("this going to be fun: %d", i),
			),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWriterTextInPlaceTags(b *testing.B) {
	ctx := context.TODO()
	writer := logw.TextLogWriter(ctx, io.Discard)

	for i := 0; i < b.N; i++ {
		_, err := writer.Write(
			[]byte(
				logw.Info.
					WithInt("attempt", i).
					WithBool("attempting", true).
					WithFloat("someFloat", 3.4).
					WithString("greeting", "Hello World").
					WithMessage("this going to be fun: %d", i),
			),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWriterJSONAllTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")
	writer := logw.JSONLogWriter(ctx, io.Discard)

	for i := 0; i < b.N; i++ {
		_, err := writer.Write(
			[]byte(
				logw.Info.
					WithInt("attempt", i).
					WithBool("attempting", true).
					WithFloat("someFloat", 3.4).
					WithString("greeting", "Hello World").
					WithMessage("this going to be fun: %d", i),
			),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWriterTextAllTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")

	writer := logw.TextLogWriter(ctx, io.Discard)
	for i := 0; i < b.N; i++ {
		_, err := writer.Write(
			[]byte(
				logw.Info.
					WithInt("attempt", i).
					WithBool("attempting", true).
					WithFloat("someFloat", 3.4).
					WithString("greeting", "Hello World").
					WithMessage("this going to be fun: %d", i),
			),
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}
