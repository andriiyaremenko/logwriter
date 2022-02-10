package logw_test

import (
	"bytes"
	"context"
	"log"
	"testing"

	logw "github.com/andriiyaremenko/logwriter"
	"github.com/stretchr/testify/suite"
)

func TestLogWriter(t *testing.T) {
	suite.Run(t, new(logWriterSuite))
}

type logWriterSuite struct {
	suite.Suite

	log log.Logger
}

func (s *logWriterSuite) TestJSONLogWriter() {
	s.Run("no_level_no_tags", func() {
		b := new(bytes.Buffer)
		test := func(result *logw.Log) {
			s.Equal(2, result.LevelCode)
			s.Equal("info", result.Level)
			s.Equal("test", result.Message)
		}

		s.log.SetOutput(logw.LogWriter(context.TODO(), b, s.getTestFormatter(test)))
		s.log.Println("test")
	})

	s.Run("with_level_no_tags", func() {
		b := new(bytes.Buffer)
		test := func(result *logw.Log) {
			s.Equal(3, result.LevelCode)
			s.Equal("warn", result.Level)
			s.Equal("test", result.Message)
		}

		s.log.SetOutput(logw.LogWriter(context.TODO(), b, s.getTestFormatter(test)))
		s.log.Println(logw.Warn, "test")
	})

	s.Run("no_level_with_tags", func() {
		b := new(bytes.Buffer)
		ctx := context.TODO()
		ctx = logw.AppendInfo(ctx, "foo", "bar")
		ctx = logw.AppendInfo(ctx, "foo", "baz")
		test := func(result *logw.Log) {
			s.Equal(2, result.LevelCode)
			s.Equal("info", result.Level)
			s.Equal("test", result.Message)
			s.ElementsMatch(
				[]logw.Tag{
					{Key: "foo", Value: "bar", Level: 2},
					{Key: "foo", Value: "baz", Level: 2},
				},
				result.Tags,
			)
		}

		s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
		s.log.Println("test")
	})

	s.Run("with_level_with_tags", func() {
		b := new(bytes.Buffer)
		ctx := context.TODO()
		ctx = logw.AppendInfo(ctx, "foo", "bar")
		ctx = logw.AppendInfo(ctx, "foo", "baz")
		test := func(result *logw.Log) {
			s.Equal(4, result.LevelCode)
			s.Equal("error", result.Level)
			s.Equal("test", result.Message)
			s.ElementsMatch(
				[]logw.Tag{
					{Key: "foo", Value: "bar", Level: 2},
					{Key: "foo", Value: "baz", Level: 2},
				},
				result.Tags,
			)
		}

		s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
		s.log.Println(logw.Error, "test")
	})

	s.Run("with_level_no_above_level_tags", func() {
		b := new(bytes.Buffer)
		ctx := context.TODO()
		ctx = logw.AppendInfo(ctx, "foo", "baz")
		ctx = logw.AppendError(ctx, "foo", "bar")
		test := func(result *logw.Log) {
			s.Equal(3, result.LevelCode)
			s.Equal("warn", result.Level)
			s.Equal("test", result.Message)
			s.ElementsMatch(
				[]logw.Tag{
					{Key: "foo", Value: "baz", Level: 2},
				},
				result.Tags,
			)
		}

		s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
		s.log.Println(logw.Warn, "test")
	})

	s.Run("with_level_in_place_tags", func() {
		b := new(bytes.Buffer)
		ctx := context.TODO()
		test := func(result *logw.Log) {
			s.Equal(2, result.LevelCode)
			s.Equal("info", result.Level)
			s.Equal("test", result.Message)
			s.ElementsMatch(
				[]logw.Tag{
					{Key: "foo", Value: true, Level: 2},
					{Key: "bar", Value: 1, Level: 2},
					{Key: "baz", Value: "test", Level: 2},
					{Key: "trace", Value: "[logwriter_test.go 137]", Level: 2},
				},
				result.Tags,
			)
		}

		s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
		s.log.Println(
			logw.Info.
				WithBool("foo", true).
				WithInt("bar", 1).
				WithString("baz", "test").
				WithRowNumber(),
			"test",
		)
	})
}

func (s *logWriterSuite) getTestFormatter(test func(*logw.Log)) logw.LogWriterOption {
	return logw.Option(2,
		func(log *logw.Log, dateF string) []byte {
			test(log)

			return logw.TextFormatter(log, dateF)
		},
		logw.NoDate,
	)
}
