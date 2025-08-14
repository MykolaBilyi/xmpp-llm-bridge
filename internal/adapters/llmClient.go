package adapters

import (
	"context"
	"xmpp-llm-bridge/internal/ports"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

type LLMClient struct {
	logger ports.Logger
	client openai.Client
}

var _ ports.LLMClient = &LLMClient{}

func NewLLMClient(config ports.Config, logger ports.Logger) (*LLMClient, error) {
	client := openai.NewClient(
		option.WithBaseURL(config.GetString("endpoint")),
		option.WithAPIKey(config.GetString("key")),
	)

	return &LLMClient{
		logger: logger,
		client: client,
	}, nil
}

func (c *LLMClient) GetChatCompletion(ctx context.Context, req ports.ChatCompletionRequest) (ports.ChatCompletionResponse, error) {
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(string(req)),
		},
		Seed:  openai.Int(0),
		Model: openai.ChatModelGPT4o,
	}

	resp, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		c.logger.Error("error creating chat completion", ports.Fields{"error": err})
		return ports.ChatCompletionResponse(""), err
	}

	return ports.ChatCompletionResponse(resp.Choices[0].Message.Content), nil
}
