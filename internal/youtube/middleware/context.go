package middleware

import "context"

type ContextKey uint8

const (
	embedKey ContextKey = iota
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
