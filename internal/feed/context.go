package feed

import (
	"context"
)

type ctxKey uint8

const formatKey ctxKey = iota

func NewFormatContext(ctx context.Context, t Format) context.Context {
	return context.WithValue(ctx, formatKey, t)
}

func FormatFromContext(ctx context.Context) (Format, bool) {
	data, ok := ctx.Value(formatKey).(Format)
	return data, ok
}
