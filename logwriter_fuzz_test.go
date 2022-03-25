package logw_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strings"
	"testing"

	logw "github.com/andriiyaremenko/logwriter"
)

func FuzzLogWriterJSONInPlaceTagsWithInt(f *testing.F) {
	ints := []int{-999, 5, 3, 0, 10, 1024, 300000, -25, 42}
	for _, tc := range ints {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, orig int) {
		b := new(bytes.Buffer)
		log := log.New(logw.JSONLogWriter(context.TODO(), b), "", log.Lmsgprefix)

		log.Println(logw.Info.WithInt("tag", orig), "some message")

		result := make(map[string]any)
		if err := json.Unmarshal(b.Bytes(), &result); err != nil {
			t.Error(err)
		}

		v, ok := result["tag"]
		if !ok {
			t.Errorf("JSONLogWriter failed for %d: no tag", orig)
			t.FailNow()
		}

		if len(v.([]any)) == 0 {
			t.Errorf("JSONLogWriter failed for %q: tag has no value", orig)
			t.FailNow()
		}

		if r := v.([]any)[0]; int(r.(float64)) != orig {
			t.Errorf("JSONLogWriter failed for %d, got %v", orig, r)
		}
	})
}

func FuzzLogWriterJSONInPlaceTagsWithFloat(f *testing.F) {
	floats := []float64{-999.2222222, 5.1, 3., 0, 10.21634, 1024.00001, 300000.26, -25.9999999299277, 42.}
	for _, tc := range floats {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, orig float64) {
		b := new(bytes.Buffer)
		log := log.New(logw.JSONLogWriter(context.TODO(), b), "", log.Lmsgprefix)

		log.Println(logw.Info.WithFloat("tag", orig), "some message")

		result := make(map[string]any)
		if err := json.Unmarshal(b.Bytes(), &result); err != nil {
			t.Error(err)
		}

		v, ok := result["tag"]
		if !ok {
			t.Errorf("JSONLogWriter failed for %f: no tag", orig)
			t.FailNow()
		}

		if len(v.([]any)) == 0 {
			t.Errorf("JSONLogWriter failed for %f: tag has no value", orig)
			t.FailNow()
		}

		if r := v.([]any)[0]; r.(float64) != orig {
			t.Errorf("JSONLogWriter failed for %f, got %v", orig, r)
		}
	})
}

func FuzzLogWriterJSONInPlaceTagsWithBool(f *testing.F) {
	bools := []bool{true, false}
	for _, tc := range bools {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, orig bool) {
		b := new(bytes.Buffer)
		log := log.New(logw.JSONLogWriter(context.TODO(), b), "", log.Lmsgprefix)

		log.Println(logw.Info.WithBool("tag", orig), "some message")

		result := make(map[string]any)
		if err := json.Unmarshal(b.Bytes(), &result); err != nil {
			t.Error(err)
		}

		v, ok := result["tag"]
		if !ok {
			t.Errorf("JSONLogWriter failed for %t: no tag", orig)
			t.FailNow()
		}

		if len(v.([]any)) == 0 {
			t.Errorf("JSONLogWriter failed for %t: tag has no value", orig)
			t.FailNow()
		}

		if r := v.([]any)[0]; r.(bool) != orig {
			t.Errorf("JSONLogWriter failed for %t, got %v", orig, r)
		}
	})
}

func FuzzLogWriterJSONInPlaceTagsWithString(f *testing.F) {
	strs := []string{"Hello, world", "true", "!12345"}
	for _, tc := range strs {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, orig string) {
		b := new(bytes.Buffer)
		log := log.New(logw.JSONLogWriter(context.TODO(), b), "", log.Lmsgprefix)

		log.Println(logw.Info.WithString("tag", orig), "some message")

		result := make(map[string]any)
		if err := json.Unmarshal(b.Bytes(), &result); err != nil {
			t.Error(err)
		}

		v, ok := result["tag"]
		if !ok {
			t.Errorf("JSONLogWriter failed for %q: no tag", orig)
			t.FailNow()
		}

		orig = strings.ReplaceAll(orig, "\t", " ")
		orig = strings.ReplaceAll(orig, "\n", " ")
		temp, _ := json.Marshal(orig)
		orig = ""
		_ = json.Unmarshal(temp, &orig)

		if len(v.([]any)) == 0 {
			t.Errorf("JSONLogWriter failed for %q: tag has no value", orig)
			t.FailNow()
		}

		if r := v.([]any)[0]; r.(string) != orig {
			t.Errorf("JSONLogWriter failed for %q, got \"%v\"", orig, r)
		}
	})
}

func FuzzLogWriterJSONWithMessage(f *testing.F) {
	strs := []string{"Hello, world", "true", "!12345"}
	for _, tc := range strs {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, orig string) {
		b := new(bytes.Buffer)
		log := log.New(logw.JSONLogWriter(context.TODO(), b), "", log.Lmsgprefix)

		log.Println(orig)

		result := make(map[string]any)
		if err := json.Unmarshal(b.Bytes(), &result); err != nil {
			t.Error(err)
			t.FailNow()
		}

		orig = strings.TrimLeft(orig, " ")
		orig = strings.TrimRight(orig, "\n")

		if orig == "" {
			return
		}

		temp, _ := json.Marshal(orig)
		orig = ""
		_ = json.Unmarshal(temp, &orig)

		if v, ok := result["message"]; !ok || v.(string) != orig {
			t.Errorf("JSONLogWriter failed for %q, got \"%v\"", orig, v)
		}
	})
}
