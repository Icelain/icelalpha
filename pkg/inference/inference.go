package inference

import (
	"context"
	"errors"
	"io"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	ollama "github.com/ollama/ollama/api"

	deepseek "github.com/cohesion-org/deepseek-go"
	deepseekconstants "github.com/cohesion-org/deepseek-go/constants"
)

type LLMClient interface {
	StreamResponse(ctx context.Context, query string) (chan string, error)
}

type ClaudeLLMClient struct {
	anthropicClient *anthropic.Client
	model           string
}

func NewClaudeLLMClient(apiKey string) *ClaudeLLMClient {

	return &ClaudeLLMClient{anthropicClient: anthropic.NewClient(option.WithAPIKey(apiKey)), model: "Claude Sonnet"}

}

// stream response to client, pass in context to manage stream cancellation
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

		close(resultChan)

	}()

	return resultChan, nil

}

// ollama based llm client works on local models

type OllamaLLMClient struct {
	ollamaclient *ollama.Client
	model        string
}

func (dl *OllamaLLMClient) StreamResponse(ctx context.Context, query string) (chan string, error) {

	to := make(chan string)
	request := &ollama.GenerateRequest{

		Model:  dl.model,
		Prompt: query,
	}

	go func() {
		dl.ollamaclient.Generate(ctx, request, func(generatedResponse ollama.GenerateResponse) error {

			to <- generatedResponse.Response
			return nil

		})

		close(to)
	}()

	return to, nil

}

func NewOllamaLLMClient(ctx context.Context, model string) (*OllamaLLMClient, error) {

	client, err := ollama.ClientFromEnvironment()
	if err != nil {

		return &OllamaLLMClient{}, nil

	}

	return &OllamaLLMClient{ollamaclient: client, model: model}, nil

}

// deepseek api based llm client

type DeepSeekLLMClient struct {
	deepseekClient *deepseek.Client
}

func NewDeepSeekLLMClient(apiKey string) *DeepSeekLLMClient {

	return &DeepSeekLLMClient{deepseekClient: deepseek.NewClient(apiKey)}

}

func (ds *DeepSeekLLMClient) StreamResponse(ctx context.Context, query string) (chan string, error) {

	request := &deepseek.StreamChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseekconstants.ChatMessageRoleUser, Content: "Just testing if the streaming feature is working or not!"},
		},
		Stream: true,
	}
	stream, err := ds.deepseekClient.CreateChatCompletionStream(ctx, request)
	if err != nil {
		return nil, err
	}

	responseChannel := make(chan string)

	go func() {

		defer stream.Close()
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			for _, choice := range response.Choices {

				responseChannel <- choice.Delta.Content

			}
		}

	}()

	return responseChannel, nil

}
