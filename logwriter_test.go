package logw_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"testing"
	"time"

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

func (s *logWriterSuite) TestOnlyMessage() {
	b := new(bytes.Buffer)
	test := func(
		level string,
		levelCode int,
		tags []logw.Tag,
		timeStamp time.Time,
		message []byte,
	) {
		s.Equal(2, levelCode)
		s.Equal("info", level)
		s.Equal("test", string(message))
	}

	s.log.SetOutput(logw.LogWriter(context.TODO(), b, s.getTestFormatter(test)))
	s.log.Println("test")
}

func (s *logWriterSuite) TestMessageAndLevel() {
	b := new(bytes.Buffer)
	test := func(
		level string,
		levelCode int,
		tags []logw.Tag,
		timeStamp time.Time,
		message []byte,
	) {
		s.Equal(3, levelCode)
		s.Equal("warn", level)
		s.Equal("test", string(message))
	}

	s.log.SetOutput(logw.LogWriter(context.TODO(), b, s.getTestFormatter(test)))
	s.log.Println(logw.Warn, "test")
}

func (s *logWriterSuite) TestMessageWithContextTags() {
	b := new(bytes.Buffer)
	ctx := context.TODO()
	ctx = logw.AppendInfo(ctx, "foo", "bar")
	ctx = logw.AppendInfo(ctx, "foo", "baz")
	test := func(
		level string,
		levelCode int,
		tags []logw.Tag,
		timeStamp time.Time,
		message []byte,
	) {
		s.Equal(2, levelCode)
		s.Equal("info", level)
		s.Equal("test", string(message))
		s.ElementsMatch(
			[]logw.Tag{
				{Key: "foo", Value: []byte("\"bar\""), Type: "json", Level: 2},
				{Key: "foo", Value: []byte("\"baz\""), Type: "json", Level: 2},
			},
			tags,
		)
	}

	s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
	s.log.Println("test")
}

func (s *logWriterSuite) TestWithMessageLevelAndContextTags() {
	b := new(bytes.Buffer)
	ctx := context.TODO()
	ctx = logw.AppendInfo(ctx, "foo", "bar")
	ctx = logw.AppendInfo(ctx, "foo", "baz")
	test := func(
		level string,
		levelCode int,
		tags []logw.Tag,
		timeStamp time.Time,
		message []byte,
	) {
		s.Equal(4, levelCode)
		s.Equal("error", level)
		s.Equal("test", string(message))
		s.ElementsMatch(
			[]logw.Tag{
				{Key: "foo", Value: []byte("\"bar\""), Type: "json", Level: 2},
				{Key: "foo", Value: []byte("\"baz\""), Type: "json", Level: 2},
			},
			tags,
		)
	}

	s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
}

func (s *logWriterSuite) TestContextTagsAppearInSameLevelOrHigher() {
	b := new(bytes.Buffer)
	ctx := context.TODO()
	ctx = logw.AppendInfo(ctx, "foo", "baz")
	ctx = logw.AppendError(ctx, "foo", "bar")
	test := func(
		level string,
		levelCode int,
		tags []logw.Tag,
		timeStamp time.Time,
		message []byte,
	) {
		s.Equal(3, levelCode)
		s.Equal("warn", level)
		s.Equal("test", string(message))
		s.ElementsMatch(
			[]logw.Tag{
				{Key: "foo", Value: []byte("\"baz\""), Type: "json", Level: 2},
			},
			tags,
		)
	}

	s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
	s.log.Println(logw.Warn, "test")
}

func (s *logWriterSuite) TestMessageWithInPlaceTags() {
	b := new(bytes.Buffer)
	ctx := context.TODO()
	test := func(
		level string,
		levelCode int,
		tags []logw.Tag,
		timeStamp time.Time,
		message []byte,
	) {
		s.Equal(2, levelCode)
		s.Equal("info", level)
		s.Equal("test", string(message))
		s.ElementsMatch(
			[]logw.Tag{
				{Key: "foo", Value: s.marshal(true), Type: "bool", Level: 2},
				{Key: "bar", Value: s.marshal(-1), Type: "int", Level: 2},
				{Key: "float", Value: s.marshal(-1.2), Type: "float64", Level: 2},
				{Key: "baz", Value: []byte("test"), Type: "string", Level: 2},
				{Key: "trace", Value: []byte("logwriter_test.go 176"), Type: "string", Level: 2},
			},
			tags,
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
}

func (s *logWriterSuite) TestFormattedMessage() {
	b := new(bytes.Buffer)
	ctx := context.TODO()
	test := func(
		level string,
		levelCode int,
		tags []logw.Tag,
		timeStamp time.Time,
		message []byte,
	) {
		s.Equal(2, levelCode)
		s.Equal("info", level)
		s.Equal("Hello World", string(message))
	}

	s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
	s.log.Println(logw.Info.WithMessage("Hello %s", "World"))
}

func (s *logWriterSuite) TestWithErrorInContextTags() {
	b := new(bytes.Buffer)
	ctx := context.TODO()
	ctx = logw.AppendInfo(ctx, "error", errors.New("some error"))

	test := func(
		level string,
		levelCode int,
		tags []logw.Tag,
		timeStamp time.Time,
		message []byte,
	) {
		s.Equal(2, levelCode)
		s.Equal("info", level)
		s.Equal("test", string(message))
		s.ElementsMatch(
			[]logw.Tag{
				{Key: "error", Value: s.marshal("some error"), Type: "json", Level: 2},
			},
			tags,
		)
	}

	s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
	s.log.Println("test")
}

func (s *logWriterSuite) TestWithError() {
	b := new(bytes.Buffer)
	test := func(
		level string,
		levelCode int,
		tags []logw.Tag,
		timeStamp time.Time,
		message []byte,
	) {
		s.Equal(2, levelCode)
		s.Equal("info", level)
		s.Equal("", string(message))
		s.ElementsMatch(
			[]logw.Tag{
				{Key: "error", Value: []byte("some error"), Type: "string", Level: 2},
			},
			tags,
		)
	}

	s.log.SetOutput(logw.LogWriter(context.TODO(), b, s.getTestFormatter(test)))
	s.log.Println(logw.Info.Error(errors.New("some error")))
}

func (s *logWriterSuite) getTestFormatter(
	test func(string, int, []logw.Tag, time.Time, []byte),
) logw.LogWriterOption {
	return logw.Option(2,
		func(
			level string,
			levelCode int,
			tags []logw.Tag,
			timeStamp time.Time,
			dateLayout string,
			message []byte,
		) []byte {
			test(level, levelCode, tags, timeStamp, message)
			return []byte{}
		},
		logw.NoDate,
	)
}

func (s *logWriterSuite) marshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		s.FailNow(err.Error())
	}

	return b
}
