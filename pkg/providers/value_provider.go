package providers

import (
	"context"

	"github.com/google/uuid"
)

type key string

type ValueProvider[T interface{}] struct {
	key key
}

func NewValueProvider[T interface{}]() *ValueProvider[T] {
	return &ValueProvider[T]{
		key: key(uuid.NewString()),
	}
}

func (vp *ValueProvider[T]) Value(ctx context.Context) *T {
	value, ok := ctx.Value(vp.key).(T)
	if ok {
		return &value
	}
	return nil
}

func (vp *ValueProvider[T]) WithValue(ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, vp.key, value)
}
