package logw_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func (s *logWriterSuite) TestLogWriter() {
	s.Run("no_level_no_tags", func() {
		b := new(bytes.Buffer)
		test := func(result *logw.Log) {
			s.Equal(2, result.LevelCode)
			s.Equal("info", result.Level)
			s.Equal("test", result.Message)
		}

		s.log.SetOutput(logw.LogWriter(context.TODO(), b, s.getTestFormatter(test)))
		s.log.Println("test")

		s.T().Log(b.String())
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

		s.T().Log(b.String())
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
					{Key: "foo", Value: []byte("\"bar\""), Level: 2},
					{Key: "foo", Value: []byte("\"baz\""), Level: 2},
				},
				result.Tags,
			)
		}

		s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
		s.log.Println("test")

		s.T().Log(b.String())
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
					{Key: "foo", Value: []byte("\"bar\""), Level: 2},
					{Key: "foo", Value: []byte("\"baz\""), Level: 2},
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
					{Key: "foo", Value: []byte("\"baz\""), Level: 2},
				},
				result.Tags,
			)
		}

		s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
		s.log.Println(logw.Warn, "test")

		s.T().Log(b.String())
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
					{Key: "foo", Value: s.marshal(true), Level: 2},
					{Key: "bar", Value: s.marshal(-1), Level: 2},
					{Key: "float", Value: s.marshal(-1.2), Level: 2},
					{Key: "baz", Value: []byte("\"test\""), Level: 2},
					{Key: "trace", Value: []byte("\"logwriter_test.go 148\""), Level: 2},
				},
				result.Tags,
			)
		}

		s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
		s.log.Println(
			logw.Info.
				WithBool("foo", true).
				WithInt("bar", -1).
				WithFloat("float", -1.2).
				WithString("baz", "test").
				WithTrace(),
			"test",
		)

		s.T().Log(b.String())
	})

	s.Run("with_message", func() {
		b := new(bytes.Buffer)
		ctx := context.TODO()
		test := func(result *logw.Log) {
			s.Equal(2, result.LevelCode)
			s.Equal("info", result.Level)
			s.Equal("Hello World", result.Message)
		}

		s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
		s.log.Println(logw.Info.WithMessage("Hello %s", "World"))

		s.T().Log(b.String())
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

func (s *logWriterSuite) marshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		s.FailNow(err.Error())
	}

	return b
}
