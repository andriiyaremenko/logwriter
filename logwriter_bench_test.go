package logw_test

import (
	"context"
	"testing"

	logw "github.com/andriiyaremenko/logwriter"
)

type w struct{}

func (w w) Write(b []byte) (int, error) {
	return 0, nil
}

func BenchmarkLogWithLogWriterJSON(b *testing.B) {
	ctx := context.TODO()
	writer := logw.JSONLogWriter(ctx, w{})

	for i := 0; i < b.N; i++ {
		if _, err := writer.Write([]byte("this going to be fun ")); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWithLogWriterText(b *testing.B) {
	ctx := context.TODO()
	writer := logw.TextLogWriter(ctx, w{})

	for i := 0; i < b.N; i++ {
		if _, err := writer.Write([]byte("this going to be fun")); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWithLogWriterJSONLevel(b *testing.B) {
	ctx := context.TODO()
	writer := logw.JSONLogWriter(ctx, w{})

	for i := 0; i < b.N; i++ {
		_, err := writer.Write([]byte(logw.Info.WithMessage("this going to be fun: %d", i)))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWithLogWriterTextLevel(b *testing.B) {
	ctx := context.TODO()
	writer := logw.TextLogWriter(ctx, w{})

	for i := 0; i < b.N; i++ {
		_, err := writer.Write([]byte(logw.Info.WithMessage("this going to be fun: %d", i)))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWithLogWriterJSONContextTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")
	writer := logw.JSONLogWriter(ctx, w{})

	for i := 0; i < b.N; i++ {
		_, err := writer.Write([]byte("this going to be fun "))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWithLogWriterTextContextTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")
	writer := logw.TextLogWriter(ctx, w{})

	for i := 0; i < b.N; i++ {
		_, err := writer.Write([]byte("this going to be fun "))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLogWithLogWriterJSONInPlaceTags(b *testing.B) {
	ctx := context.TODO()
	writer := logw.JSONLogWriter(ctx, w{})

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

func BenchmarkLogWithLogWriterTextInPlaceTags(b *testing.B) {
	ctx := context.TODO()
	writer := logw.TextLogWriter(ctx, w{})

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

func BenchmarkLogWithLogWriterJSONAllTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")
	writer := logw.JSONLogWriter(ctx, w{})

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

func BenchmarkLogWithLogWriterTextAllTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")

	writer := logw.TextLogWriter(ctx, w{})
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
