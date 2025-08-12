package providers_test

import (
	"__template__/internal/adapters"
	"__template__/internal/providers"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type contextKey string

func TestLoggerProvider_WhenWithLoggerCalled_ReturnsChildContext(t *testing.T) {
	// Arrange
	parentTestKey := contextKey("parent key")
	parentTestValue := "test value"
	ctx := context.WithValue(context.Background(), parentTestKey, parentTestValue)
	defaultLogger := adapters.NewNoopLogger()
	provider := providers.NewLoggerProvider(defaultLogger)

	// Act
	childContext := provider.WithLogger(ctx, defaultLogger)

	// Assert
	assert.NotNil(t, childContext, "child context should not be nil")
	assert.NotSame(t, ctx, childContext, "child context should not be the same as the parent context")
	assert.Equal(t, childContext.Value(parentTestKey), parentTestValue, "child context contains same value as parent context")
}

func TestLoggerProvider_WhenLoggerInContext_ReturnsTheLogger(t *testing.T) {
	// Arrange
	defaultLogger := adapters.NewNoopLogger()
	provider := providers.NewLoggerProvider(defaultLogger)
	contextLogger := adapters.NewNoopLogger()
	childContext := provider.WithLogger(context.Background(), contextLogger)

	// Act
	logger := provider.Value(childContext)

	// Assert
	assert.NotNil(t, logger, "logger should not be nil")
	assert.NotSame(t, defaultLogger, logger, "logger should not be the default logger")
	assert.Same(t, contextLogger, logger, "logger should be the same as the one in the context")
}

func TestLoggerProvider_WhenNoLoggerInContext_ReturnsDefaultLogger(t *testing.T) {
	// Arrange
	defaultLogger := adapters.NewNoopLogger()
	provider := providers.NewLoggerProvider(defaultLogger)

	// Act
	logger := provider.Value(context.Background())

	// Assert
	assert.NotNil(t, logger, "logger should not be nil")
	assert.Same(t, defaultLogger, logger, "logger should be the default logger")
}
