package ports

import "context"

type Fields = map[string]interface{}

//go:generate mockgen -destination ../adapters/mocks/logger.go -package adapter_mocks . Logger

type Logger interface {
	Trace(string, ...Fields)
	Debug(string, ...Fields)
	Info(string, ...Fields)
	Warn(string, ...Fields)
	Error(string, ...Fields)

	WithFields(Fields) Logger
	WithContext(context.Context) Logger
}

type LoggerProvider interface {
	WithLogger(context.Context, Logger) context.Context
	Value(context.Context) Logger
}
