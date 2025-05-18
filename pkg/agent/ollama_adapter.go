package agent

import (
	"context"

	"github.com/ollama/ollama/api"
)

type ollamaProvider struct {
	client *api.Client
	model  string
}

// NewOllamaProvider builds a Provider backed by the Ollama client.
func NewOllamaProvider(client *api.Client, model string) Provider {
	return &ollamaProvider{client: client, model: model}
}

func (o *ollamaProvider) Chat(ctx context.Context, messages []MessageParam) (Message, error) {
	var chatMsgs []api.ChatMessage
	for _, m := range messages {
		chatMsgs = append(chatMsgs, api.ChatMessage{
			Role:    string(m.Role),
			Content: m.Content,
		})
	}
	req := api.ChatRequest{
		Model:    o.model,
		Messages: chatMsgs,
	}
	return o.client.Chat(ctx, req)
}
