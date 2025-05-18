package agent

import "context"

// Provider is the minimal chat-completion API that Agent needs.
type Provider interface {
	Chat(ctx context.Context, messages []MessageParam) (Message, error)
}
