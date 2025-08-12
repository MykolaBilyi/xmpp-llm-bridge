package providers

import (
	"context"

	"xmpp-llm-bridge/internal/ports"

	"github.com/google/uuid"
)

type key string

type LoggerProvider struct {
	key           key
	defaultLogger ports.Logger
}

func NewLoggerProvider(defaultLogger ports.Logger) *LoggerProvider {
	return &LoggerProvider{
		key:           key(uuid.NewString()),
		defaultLogger: defaultLogger,
	}
}

var _ ports.LoggerProvider = &LoggerProvider{}

func (vp *LoggerProvider) Value(ctx context.Context) ports.Logger {
	value, ok := ctx.Value(vp.key).(ports.Logger)
	if ok {
		return value
	}
	return vp.defaultLogger
}

func (vp *LoggerProvider) WithLogger(ctx context.Context, value ports.Logger) context.Context {
	return context.WithValue(ctx, vp.key, value)
}
