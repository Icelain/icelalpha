package inference

import (
	"context"

	anthropic "github.com/liushuangls/go-anthropic"
)

type LLMClient interface {
	StreamResponse(ctx context.Context, query string) (<-chan string, error)
}

type ClaudeLLMClient struct {
	anthropicClient *anthropic.Client
	model           string
}

func NewClaudeLLMClient(apiKey string) *ClaudeLLMClient {

	return &ClaudeLLMClient{anthropicClient: anthropic.NewClient(apiKey), model: "Claude Sonnet"}

}

func (c *ClaudeLLMClient) StreamResponse(ctx context.Context, query string) (<-chan string, error) {

	request := anthropic.MessagesStreamRequest{

		MessagesRequest: anthropic.MessagesRequest{

			Model: anthropic.ModelClaude3Sonnet20240229,
			Messages: []anthropic.Message{

				{Role: anthropic.RoleUser, Content: []anthropic.MessageContent{

					{Text: &query},
				}},
			},
			Stream: true,
		},
	}

	response, err := c.anthropicClient.CreateMessagesStream(ctx, request)
	if err != nil {

	}

}
