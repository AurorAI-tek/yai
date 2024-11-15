package ai

import (
	"context"
)

// ProviderType represents different AI provider types
type ProviderType string

const (
	OpenAIProvider ProviderType = "openai"
	GroqProvider   ProviderType = "groq"
	LlamaProvider  ProviderType = "llama"
)

// Provider interface defines methods that any AI provider must implement
type Provider interface {
	CreateChatCompletion(ctx context.Context, messages []Message) (*CompletionResponse, error)
	CreateChatCompletionStream(ctx context.Context, messages []Message) (Stream, error)
}

// Message represents a chat message
type Message struct {
	Role    string
	Content string
}

// CompletionResponse represents a response from the AI provider
type CompletionResponse struct {
	Content string
}

// Stream represents a streaming response interface
type Stream interface {
	Recv() (*CompletionResponse, error)
	Close()
}
