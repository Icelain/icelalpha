package inference

import (
	"context"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type LLMClient interface {
	StreamResponse(ctx context.Context, query string) (<-chan string, error)
}

type ClaudeLLMClient struct {
	anthropicClient *anthropic.Client
	model           string
}

func NewClaudeLLMClient(apiKey string) *ClaudeLLMClient {

	return &ClaudeLLMClient{anthropicClient: anthropic.NewClient(option.WithAPIKey(apiKey)), model: "Claude Sonnet"}

}

func (cl *ClaudeLLMClient) StreamResponse(ctx context.Context, query string) (chan string, error) {

	stream := cl.anthropicClient.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
		MaxTokens: anthropic.Int(1024),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(query)),
		}),
	})

	if err := stream.Err(); err != nil {

		return nil, err

	}

	message := anthropic.Message{}
	resultChan := make(chan string)

	go func() {

		for stream.Next() {
			event := stream.Current()
			message.Accumulate(event)

			switch delta := event.Delta.(type) {
			case anthropic.ContentBlockDeltaEventDelta:
				if delta.Text != "" {

					// we have a stream response from claude
					resultChan <- delta.Text

				}
			}
		}

	}()

	return resultChan, nil

}
