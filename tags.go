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

func AppendDebug(ctx context.Context, tag string, value interface{}) context.Context {
	return AppendTag(ctx, levelDebug, tag, value)
}

func AppendInfo(ctx context.Context, tag string, value interface{}) context.Context {
	return AppendTag(ctx, levelInfo, tag, value)
}

func AppendWarn(ctx context.Context, tag string, value interface{}) context.Context {
	return AppendTag(ctx, levelWarn, tag, value)
}

func AppendError(ctx context.Context, tag string, value interface{}) context.Context {
	return AppendTag(ctx, levelError, tag, value)
}

func AppendFatal(ctx context.Context, tag string, value interface{}) context.Context {
	return AppendTag(ctx, levelFatal, tag, value)
}

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
