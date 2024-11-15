package groq

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ekkinox/yai/ai"
)

const (
	groqAPIEndpoint = "https://api.groq.com/v1/chat/completions"
)

type Provider struct {
	apiKey string
	model  string
}

func NewProvider(apiKey, model string) *Provider {
	return &Provider{
		apiKey: apiKey,
		model:  model,
	}
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqRequest struct {
	Model    string        `json:"model"`
	Messages []groqMessage `json:"messages"`
	Stream   bool         `json:"stream"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (p *Provider) CreateChatCompletion(ctx context.Context, messages []ai.Message) (*ai.CompletionResponse, error) {
	groqMessages := make([]groqMessage, len(messages))
	for i, msg := range messages {
		groqMessages[i] = groqMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	reqBody, err := json.Marshal(groqRequest{
		Model:    p.model,
		Messages: groqMessages,
		Stream:   false,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", groqAPIEndpoint, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var groqResp groqResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return nil, err
	}

	if len(groqResp.Choices) == 0 {
		return nil, fmt.Errorf("no completion choices returned")
	}

	return &ai.CompletionResponse{
		Content: groqResp.Choices[0].Message.Content,
	}, nil
}

func (p *Provider) CreateChatCompletionStream(ctx context.Context, messages []ai.Message) (ai.Stream, error) {
	groqMessages := make([]groqMessage, len(messages))
	for i, msg := range messages {
		groqMessages[i] = groqMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	reqBody, err := json.Marshal(groqRequest{
		Model:    p.model,
		Messages: groqMessages,
		Stream:   true,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", groqAPIEndpoint, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return &groqStream{
		reader: resp.Body,
	}, nil
}

type groqStream struct {
	reader io.ReadCloser
}

func (s *groqStream) Recv() (*ai.CompletionResponse, error) {
	var buf [1024]byte
	n, err := s.reader.Read(buf[:])
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("error reading stream: %v", err)
	}

	data := buf[:n]
	if len(data) == 0 {
		return nil, io.EOF
	}

	// Parse SSE data
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			content := strings.TrimPrefix(line, "data: ")
			return &ai.CompletionResponse{Content: content}, nil
		}
	}

	return nil, nil
}

func (s *groqStream) Close() {
	s.reader.Close()
}
