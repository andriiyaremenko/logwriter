package logw_test

import (
	"bytes"
	"context"
	"io"
	"log"
	"testing"

	logw "github.com/andriiyaremenko/logwriter"
)

func BenchmarkLogToCompare(b *testing.B) {
	var writer io.Writer = new(bytes.Buffer)
	log := getLog(writer)

	for i := 0; i < b.N; i++ {
		log.Println("this going to be fun ", i)
	}
}

func BenchmarkLogWithLogWriterJSON(b *testing.B) {
	ctx := context.TODO()
	writer := logw.JSONLogWriter(ctx, new(bytes.Buffer))
	log := getLog(writer)

	for i := 0; i < b.N; i++ {
		log.Println("this going to be fun", i)
	}
}

func BenchmarkLogWithLogWriterText(b *testing.B) {
	ctx := context.TODO()
	writer := logw.TextLogWriter(ctx, new(bytes.Buffer))
	log := getLog(writer)

	for i := 0; i < b.N; i++ {
		log.Println("this going to be fun", i)
	}
}

func BenchmarkLogLevelToCompare(b *testing.B) {
	var writer io.Writer = new(bytes.Buffer)
	log := getLog(writer)

	for i := 0; i < b.N; i++ {
		log.Println("this going to be fun ", i)
	}
}

func BenchmarkLogWithLogWriterJSONLevel(b *testing.B) {
	ctx := context.TODO()
	writer := logw.JSONLogWriter(ctx, new(bytes.Buffer))
	log := getLog(writer)

	for i := 0; i < b.N; i++ {
		log.Println(
			logw.Info.WithMessage("this going to be fun: %d", i),
		)
	}
}

func BenchmarkLogWithLogWriterTextLevel(b *testing.B) {
	ctx := context.TODO()
	writer := logw.TextLogWriter(ctx, new(bytes.Buffer))
	log := getLog(writer)

	for i := 0; i < b.N; i++ {
		log.Println(
			logw.Info.WithMessage("this going to be fun: %d", i),
		)
	}
}

func BenchmarkLogContextTagsToCompare(b *testing.B) {
	var writer io.Writer = new(bytes.Buffer)
	log := getLog(writer)

	for i := 0; i < b.N; i++ {
		log.Println(
			"tag1", true,
			"tag2", 42,
			"tag3", "Hello World",
			"this going to be fun ", i,
		)
	}
}

func BenchmarkLogWithLogWriterJSONContextTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")
	writer := logw.JSONLogWriter(ctx, new(bytes.Buffer))
	log := getLog(writer)

	for i := 0; i < b.N; i++ {
		log.Println("this going to be fun ", i)
	}
}

func BenchmarkLogWithLogWriterTextContextTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")
	writer := logw.TextLogWriter(ctx, new(bytes.Buffer))
	log := getLog(writer)

	for i := 0; i < b.N; i++ {
		log.Println("this going to be fun ", i)
	}
}

func BenchmarkLogInPlaceTagsToCompare(b *testing.B) {
	var writer io.Writer = new(bytes.Buffer)
	log := getLog(writer)

	for i := 0; i < b.N; i++ {
		log.Println(
			"attempt", i,
			"attempting", true,
			"someFloat", 3.4,
			"greeting", "Hello World",
			"this going to be fun ", i,
		)
	}
}

func BenchmarkLogWithLogWriterJSONInPlaceTags(b *testing.B) {
	ctx := context.TODO()
	writer := logw.JSONLogWriter(ctx, new(bytes.Buffer))
	log := getLog(writer)

	for i := 0; i < b.N; i++ {
		log.Println(
			logw.Info.
				WithInt("attempt", i).
				WithBool("attempting", true).
				WithFloat("someFloat", 3.4).
				WithString("greeting", "Hello World").
				WithMessage("this going to be fun: %d", i),
		)
	}
}

func BenchmarkLogWithLogWriterTextInPlaceTags(b *testing.B) {
	ctx := context.TODO()
	writer := logw.TextLogWriter(ctx, new(bytes.Buffer))
	log := getLog(writer)

	for i := 0; i < b.N; i++ {
		log.Println(
			logw.Info.
				WithInt("attempt", i).
				WithBool("attempting", true).
				WithFloat("someFloat", 3.4).
				WithString("greeting", "Hello World").
				WithMessage("this going to be fun: %d", i),
		)
	}
}

func BenchmarkLogAllTagsToCompare(b *testing.B) {
	var writer io.Writer = new(bytes.Buffer)
	log := getLog(writer)
	for i := 0; i < b.N; i++ {
		log.Println(
			"tag1", true,
			"tag2", 42,
			"tag3", "Hello World",
			"attempt", i,
			"attempting", true,
			"someFloat", 3.4,
			"greeting", "Hello World",
			"this going to be fun ", i,
		)
	}
}

func BenchmarkLogWithLogWriterJSONAllTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")
	writer := logw.JSONLogWriter(ctx, new(bytes.Buffer))
	log := getLog(writer)

	for i := 0; i < b.N; i++ {
		log.Println(
			logw.Info.
				WithInt("attempt", i).
				WithBool("attempting", true).
				WithFloat("someFloat", 3.4).
				WithString("greeting", "Hello World").
				WithMessage("this going to be fun: %d", i),
		)
	}
}

func BenchmarkLogWithLogWriterTextAllTags(b *testing.B) {
	ctx := context.TODO()
	logw.AppendInfo(ctx, "tag1", true)
	logw.AppendInfo(ctx, "tag2", 42)
	logw.AppendInfo(ctx, "tag3", "Hello World")

	writer := logw.TextLogWriter(ctx, new(bytes.Buffer))
	log := getLog(writer)
	for i := 0; i < b.N; i++ {
		log.Println(
			logw.Info.
				WithInt("attempt", i).
				WithBool("attempting", true).
				WithFloat("someFloat", 3.4).
				WithString("greeting", "Hello World").
				WithMessage("this going to be fun: %d", i),
		)
	}
}

func getLog(w io.Writer) *log.Logger {
	return log.New(w, "", log.Lmsgprefix)
}
