package agent

import (
	"context"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

type anthropicProvider struct {
	client *anthropic.Client
	opts   anthropic.MessageNewParams
}

// NewAnthropicProvider builds a Provider backed by the Anthropic client.
func NewAnthropicProvider(client *anthropic.Client, opts anthropic.MessageNewParams) Provider {
	return &anthropicProvider{client: client, opts: opts}
}

func (a *anthropicProvider) Chat(ctx context.Context, messages []MessageParam) (Message, error) {
	return a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     a.opts.Model,
		MaxTokens: a.opts.MaxTokens,
		Messages:  messages,
		Tools:     a.opts.Tools,
	})
}
