package providers

import (
	"context"

	"xmpp-llm-bridge/internal/entities"
	"xmpp-llm-bridge/pkg/providers"
)

//go:generate mockgen -destination mocks/request_id_provider.go -package provider_mocks . RequestIdProvider

type RequestIdProvider interface {
	Value(ctx context.Context) *entities.RequestId
	WithValue(ctx context.Context, value entities.RequestId) context.Context
}

func NewRequestIdProvider() RequestIdProvider {
	return providers.NewValueProvider[entities.RequestId]()
}
