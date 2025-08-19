package adapters

import (
	"context"
	"xmpp-llm-bridge/internal/ports"
	"xmpp-llm-bridge/internal/providers"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

type OpenAIClient struct {
	config         ports.Config
	loggerProvider *providers.LoggerProvider
	client         openai.Client
}

var _ ports.LLMService = &OpenAIClient{}

func NewOpenAIClient(ctx context.Context, config ports.Config, loggerProvider *providers.LoggerProvider) (*OpenAIClient, error) {
	config.SetDefault("api.endpoint", "https://api.openai.com/v1")
	config.SetDefault("api.model", "gpt-3.5-turbo")
	config.SetDefault("maxTokens", 800)
	config.SetDefault("temperature", 0.7)
	config.SetDefault("seed", 0)

	client := openai.NewClient(
		option.WithBaseURL(config.GetString("api.endpoint")),
		option.WithAPIKey(config.GetString("api.key")),
	)

	return &OpenAIClient{
		config:         config,
		loggerProvider: loggerProvider,
		client:         client,
	}, nil
}

func (c *OpenAIClient) GetChatCompletion(ctx context.Context, req ports.ChatCompletionRequest) (ports.ChatCompletionResponse, error) {
	logger := c.loggerProvider.Value(ctx)
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(c.config.GetString("systemPrompt")),
			openai.UserMessage(string(req)),
		},
		Seed:        openai.Int(int64(c.config.GetInt("seed"))),
		Model:       c.config.GetString("api.model"),
		MaxTokens:   openai.Int(int64(c.config.GetInt("maxTokens"))),
		Temperature: openai.Float(c.config.GetFloat64("temperature")),
	}

	resp, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		logger.Error("error creating chat completion", ports.Fields{"error": err})
		return ports.ChatCompletionResponse(""), err
	}

	return ports.ChatCompletionResponse(resp.Choices[0].Message.Content), nil
}
