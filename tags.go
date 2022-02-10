package logw

import (
	"context"
)

type Tag struct {
	Key   string
	Value interface{}
	Level int
}

type key int

var logwriterKey key

// Addends Tag to context, that will be logged with Debug level
func AppendDebug(ctx context.Context, tag string, value interface{}) context.Context {
	return AppendTag(ctx, LevelDebug, tag, value)
}

// Addends Tag to context, that will be logged with Info level
func AppendInfo(ctx context.Context, tag string, value interface{}) context.Context {
	return AppendTag(ctx, LevelInfo, tag, value)
}

// Addends Tag to context, that will be logged with Warn level
func AppendWarn(ctx context.Context, tag string, value interface{}) context.Context {
	return AppendTag(ctx, LevelWarn, tag, value)
}

// Addends Tag to context, that will be logged with Error level
func AppendError(ctx context.Context, tag string, value interface{}) context.Context {
	return AppendTag(ctx, LevelError, tag, value)
}

// Addends Tag to context, that will be logged with Fatal level
func AppendFatal(ctx context.Context, tag string, value interface{}) context.Context {
	return AppendTag(ctx, LevelFatal, tag, value)
}

// Addends Tag to context, that will be logged with provided level
func AppendTag(ctx context.Context, level int, tag string, value interface{}) context.Context {
	newTag := Tag{
		Key:   tag,
		Value: value,
		Level: level,
	}

	tags, ok := ctx.Value(logwriterKey).([]Tag)
	if !ok {
		return context.WithValue(ctx, logwriterKey, []Tag{newTag})
	}

	for _, oldTag := range tags {
		if oldTag == newTag {
			return context.WithValue(ctx, logwriterKey, tags)
		}
	}

	return context.WithValue(ctx, logwriterKey, append(tags, newTag))
}

func getTags(ctx context.Context, level int) []Tag {
	result := make([]Tag, 0, 1)
	tags, ok := ctx.Value(logwriterKey).([]Tag)
	if !ok {
		return result
	}

	for _, tagValue := range tags {
		if tagValue.Level <= level {
			result = append(result, tagValue)
		}
	}

	return result
}
