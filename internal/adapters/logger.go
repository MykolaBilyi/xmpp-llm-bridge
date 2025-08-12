package adapters

import (
	"__template__/internal/ports"
	"context"

	zaplogfmt "github.com/sykesm/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	zapadapter "logur.dev/adapter/zap"
	"logur.dev/logur"
)

// Logger wraps a logur logger and exposes it under a custom interface.
type Logger struct {
	logur.LoggerFacade

	extractor ContextExtractor
}

var _ ports.Logger = &Logger{}

// ContextExtractor extracts log fields from a context.
type ContextExtractor func(ctx context.Context) ports.Fields

func init() {
	zap.RegisterEncoder("logfmt", func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return zaplogfmt.NewEncoder(cfg), nil
	})
}

func makeZapConfig(config ports.Config) zap.Config {
	atomicLogLevel, err := zap.ParseAtomicLevel(config.GetString("level"))
	if err != nil {
		atomicLogLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = config.GetString("timeKey")
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	return zap.Config{
		Level:             atomicLogLevel,
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          config.GetString("format"),
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}
}

// NewLogger returns a new Logger instance.
func NewLogger(config ports.Config) (*Logger, error) {
	config.SetDefault("level", "info")
	config.SetDefault("timeKey", "time")
	config.SetDefault("format", "json")

	zapLogger, err := makeZapConfig(config).Build()
	if err != nil {
		return nil, err
	}

	defer zapLogger.Sync()

	return &Logger{
		LoggerFacade: zapadapter.New(zapLogger),
	}, nil
}

// NewContextAwareLogger returns a new Logger instance that can extract information from a context.
func NewContextAwareLogger(logger logur.LoggerFacade, extractor ContextExtractor) *Logger {
	return &Logger{
		LoggerFacade: logur.WithContextExtractor(logger, logur.ContextExtractor(extractor)),
		extractor:    extractor,
	}
}

// WithFields annotates a logger with key-value pairs.
func (l *Logger) WithFields(fields ports.Fields) ports.Logger {
	return &Logger{
		LoggerFacade: logur.WithFields(l.LoggerFacade, fields),
		extractor:    l.extractor,
	}
}

// WithContext annotates a logger with a context.
func (l *Logger) WithContext(ctx context.Context) ports.Logger {
	if l.extractor == nil {
		return l
	}

	return l.WithFields(l.extractor(ctx))
}

func NewNoopLogger() *Logger {
	return &Logger{
		LoggerFacade: &logur.NoopLogger{},
	}
}
