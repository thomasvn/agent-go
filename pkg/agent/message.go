package agent

import anthropic "github.com/anthropics/anthropic-sdk-go"

// MessageParam is alias for the Anthropic SDK's chat message parameter.
type MessageParam = anthropic.MessageParam

// Message is alias for the Anthropic SDK's chat completion response.
type Message = *anthropic.Message
