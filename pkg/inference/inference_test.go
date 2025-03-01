package inference

import (
	"context"
	"testing"
)

func TestClaudeStreamResponse(t *testing.T) {

	// implement later
	// pass by default
}

func TestOllamaDeepseekStreamResponse(t *testing.T) {

	ollamaClient, err := NewOllamaLLMClient(context.Background(), "deepseek-r1:14b")
	if err != nil {

		t.Fatal("error creating ollama client", err)

	}

	respchan, err := ollamaClient.StreamResponse(context.Background(), "what is the derivative of ln(x)? Only and only give the answer to the question and do not attempt to say anything else")
	if err != nil {

		t.Fatal("error creating stream response from ollama client with model: ", ollamaClient.model)

	}

	for msg := range respchan {

		t.Log("output tokens: ", msg)

	}

}

func TestDeepSeekAPIStreamResponse(t *testing.T) {

	// implement later
}
