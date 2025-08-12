package providers_test

import (
	"__template__/pkg/providers"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type contextKey string

func TestValueProvider_WhenWithValueCalled_ReturnsChildContext(t *testing.T) {
	// Arrange
	parentTestKey := contextKey("parent key")
	parentTestValue := "test value"
	ctx := context.WithValue(context.Background(), parentTestKey, parentTestValue)
	provider := providers.NewValueProvider[string]()

	// Act
	childContext := provider.WithValue(ctx, "some value")

	// Assert
	assert.NotNilf(t, childContext, "child context should not be nil")
	assert.NotEqualf(t, ctx, childContext, "child context should not be the same as the parent context")
	assert.Equalf(t, childContext.Value(parentTestKey), parentTestValue, "parent value should be the same in the child context")
}

func TestValueProvider_WhenValueInContext_ReturnsTheValue(t *testing.T) {
	// Arrange
	ctx := context.Background()
	value := "test value"
	provider := providers.NewValueProvider[string]()
	childContext := provider.WithValue(ctx, value)

	// Act
	result := provider.Value(childContext)

	// Assert
	assert.NotNilf(t, result, "result should not be nil")
	assert.Equalf(t, value, *result, "result should be the same as the value")
}

func TestValueProvider_WhenNoValueInContext_ReturnsNil(t *testing.T) {
	// Arrange
	ctx := context.Background()
	provider := providers.NewValueProvider[string]()

	// Act
	result := provider.Value(ctx)

	// Assert
	assert.Nil(t, result, "result should be nil")
}
