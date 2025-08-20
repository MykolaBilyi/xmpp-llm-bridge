package providers

import (
	"xmpp-llm-bridge/internal/entities"
	"xmpp-llm-bridge/pkg/providers"
)

//go:generate mockgen -destination mocks/request_id_provider.go -package provider_mocks . RequestIdProvider

type RequestIdProvider = providers.ValueProvider[entities.RequestId]

func NewRequestIdProvider() *RequestIdProvider {
	return providers.NewValueProvider[entities.RequestId]()
}
