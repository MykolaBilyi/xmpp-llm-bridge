package ports

import (
	"context"
)

type LLMService interface {
	GetChatCompletion(ctx context.Context, req ChatCompletionRequest) (ChatCompletionResponse, error)
}

type ChatCompletionResponse string
type ChatCompletionRequest string
