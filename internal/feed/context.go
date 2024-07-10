package feed

import (
	"context"

	"github.com/eduncan911/podcast"
	"github.com/gorilla/feeds"
)

type ctxKey uint8

const (
	formatKey ctxKey = iota
	feedKey
)

type Feed interface {
	*feeds.Feed | *podcast.Podcast
}

func NewContext[T Feed](ctx context.Context, data T) context.Context {
	return context.WithValue(ctx, feedKey, data)
}

func FromContext[T Feed | any](ctx context.Context) (T, bool) {
	data, ok := ctx.Value(feedKey).(T)
	return data, ok
}

func NewFormatContext(ctx context.Context, t Format) context.Context {
	return context.WithValue(ctx, formatKey, t)
}

func FormatFromContext(ctx context.Context) (Format, bool) {
	data, ok := ctx.Value(formatKey).(Format)
	return data, ok
}
