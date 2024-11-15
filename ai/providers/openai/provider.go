package openai

import (
	"context"
	"io"

	"github.com/ekkinox/yai/ai"
	gopenai "github.com/sashabaranov/go-openai"
)

type Provider struct {
	client *gopenai.Client
	model  string
}

func NewProvider(client *gopenai.Client, model string) *Provider {
	return &Provider{
		client: client,
		model:  model,
	}
}

func (p *Provider) CreateChatCompletion(ctx context.Context, messages []ai.Message) (*ai.CompletionResponse, error) {
	openaiMessages := make([]gopenai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = gopenai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	resp, err := p.client.CreateChatCompletion(
		ctx,
		gopenai.ChatCompletionRequest{
			Model:    p.model,
			Messages: openaiMessages,
		},
	)
	if err != nil {
		return nil, err
	}

	return &ai.CompletionResponse{
		Content: resp.Choices[0].Message.Content,
	}, nil
}

func (p *Provider) CreateChatCompletionStream(ctx context.Context, messages []ai.Message) (ai.Stream, error) {
	openaiMessages := make([]gopenai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = gopenai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	stream, err := p.client.CreateChatCompletionStream(
		ctx,
		gopenai.ChatCompletionRequest{
			Model:    p.model,
			Messages: openaiMessages,
			Stream:   true,
		},
	)
	if err != nil {
		return nil, err
	}

	return &openaiStream{stream: stream}, nil
}

type openaiStream struct {
	stream *gopenai.ChatCompletionStream
}

func (s *openaiStream) Recv() (*ai.CompletionResponse, error) {
	resp, err := s.stream.Recv()
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, err
	}

	return &ai.CompletionResponse{
		Content: resp.Choices[0].Delta.Content,
	}, nil
}

func (s *openaiStream) Close() {
	s.stream.Close()
}
