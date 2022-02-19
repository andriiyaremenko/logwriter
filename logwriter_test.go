package logw_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func (s *logWriterSuite) TestLogWriter() {
	s.Run("no_level_no_tags", func() {
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

		s.T().Log(b.String())
	})

	s.Run("with_level_no_tags", func() {
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

		s.T().Log(b.String())
	})

	s.Run("no_level_with_tags", func() {
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
					{Key: "foo", Value: []byte("\"bar\""), Level: 2},
					{Key: "foo", Value: []byte("\"baz\""), Level: 2},
				},
				tags,
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
					{Key: "foo", Value: []byte("\"bar\""), Level: 2},
					{Key: "foo", Value: []byte("\"baz\""), Level: 2},
				},
				tags,
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
					{Key: "foo", Value: []byte("\"baz\""), Level: 2},
				},
				tags,
			)
		}

		s.log.SetOutput(logw.LogWriter(ctx, b, s.getTestFormatter(test)))
		s.log.Println(logw.Warn, "test")

		s.T().Log(b.String())
	})

	s.Run("with_level_in_place_tags", func() {
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
					{Key: "foo", Value: s.marshal(true), Level: 2},
					{Key: "bar", Value: s.marshal(-1), Level: 2},
					{Key: "float", Value: s.marshal(-1.2), Level: 2},
					{Key: "baz", Value: []byte("\"test\""), Level: 2},
					{Key: "trace", Value: []byte("\"logwriter_test.go 185\""), Level: 2},
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

		s.T().Log(b.String())
	})

	s.Run("with_message", func() {
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

		s.T().Log(b.String())
	})
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

			return logw.TextFormatter(
				level,
				levelCode,
				tags,
				timeStamp,
				dateLayout,
				message,
			)
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
