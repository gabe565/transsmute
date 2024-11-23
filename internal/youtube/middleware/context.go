package middleware

import "context"

type ContextKey uint8

const (
	embedKey ContextKey = iota
	limitKey
)

func NewEmbedContext(ctx context.Context, v bool) context.Context {
	return context.WithValue(ctx, embedKey, v)
}

func EmbedFromContext(ctx context.Context) bool {
	v, ok := ctx.Value(embedKey).(bool)
	if !ok {
		return true
	}
	return v
}

func NewLimitContext(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, limitKey, v)
}

func LimitFromContext(ctx context.Context) int {
	v, ok := ctx.Value(limitKey).(int)
	if !ok {
		return DefaultLimit
	}
	return v
}
