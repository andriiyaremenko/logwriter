package logw_test

import (
	"encoding/json"
	"testing"
	"time"

	logw "github.com/andriiyaremenko/logwriter"
	"github.com/stretchr/testify/suite"
)

func TestJSONFormatter(t *testing.T) {
	suite.Run(t, new(jsonFormatterSuite))
}

type jsonFormatterSuite struct {
	suite.Suite
}

func (s *jsonFormatterSuite) TestProducesValidJSON() {
	tags := []logw.Tag{
		{Key: "foo", Value: s.marshal(true), Level: 2},
		{Key: "bar", Value: s.marshal(-1), Level: 2},
		{Key: "float", Value: s.marshal(-1.2), Level: 2},
		{Key: "baz", Value: []byte("\"test\""), Level: 2},
		{Key: "trace", Value: []byte("\"logwriter_test.go 29\""), Level: 2},
		{Key: "some_slice", Value: s.marshal([]interface{}{1, true, "test", 1.18}), Level: 2},
	}
	b := logw.JSONFormatter("info", 2, tags, time.Now(), time.RFC3339, "test json output")

	s.T().Log(string(b))

	result := make(map[string]interface{})
	err := json.Unmarshal(b, &result)

	s.NoError(err)

	for _, tag := range tags {
		if _, ok := result[tag.Key]; !ok {
			s.Fail("tag %q was dropped", tag.Key)
		}
	}
}

func (s *jsonFormatterSuite) marshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		s.FailNow(err.Error())
	}

	return b
}
