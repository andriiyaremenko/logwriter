package logw_test

import (
	"bytes"
	"context"
	"io"
	"log"
	"testing"

	logw "github.com/andriiyaremenko/logwriter"
)

func BenchmarkLog(b *testing.B) {
	var writer io.Writer = new(bytes.Buffer)
	log := getLog(writer)
	for i := 0; i < b.N; i++ {
		log.Printf("this going to be fun: %d", i)
	}
}

func BenchmarkLogWithLogWriter(b *testing.B) {
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

func getLog(w io.Writer) *log.Logger {
	return log.New(w, "", log.Lmsgprefix)
}
