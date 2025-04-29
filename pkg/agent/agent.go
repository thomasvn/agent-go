package agent

import (
	"agent/pkg/tool"
	"context"
	"encoding/json"
	"fmt"

	"agent/pkg/mcp"

	"github.com/anthropics/anthropic-sdk-go"
)

type Agent struct {
	client         *anthropic.Client
	getUserMessage func() (string, bool)
	tools          []tool.ToolDefinition
	mcpManager     *mcp.Manager
}

func NewAgent(client *anthropic.Client, getUserMessage func() (string, bool), tools []tool.ToolDefinition, mcpManager *mcp.Manager) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
		mcpManager:     mcpManager,
	}
}

func (a *Agent) Run(ctx context.Context) error {
	conversation := []anthropic.MessageParam{}

	fmt.Println("\nChat with Claude (use 'ctrl-c' to quit)")

	readUserInput := true
	for {
		if readUserInput {
			fmt.Print("\u001b[94mYou\u001b[0m: ")
			userInput, ok := a.getUserMessage()
			if !ok {
				break
			}

			userMessage := anthropic.NewUserMessage(anthropic.NewTextBlock(userInput))
			conversation = append(conversation, userMessage)
		}

		message, err := a.runInference(ctx, conversation)
		if err != nil {
			return err
		}
		conversation = append(conversation, message.ToParam())

		toolResults := []anthropic.ContentBlockParamUnion{}
		for _, content := range message.Content {
			switch content.Type {
			case "text":
				fmt.Printf("\u001b[93mClaude\u001b[0m: %s\n", content.Text)
			case "tool_use":
				if a.isLocalTool(content.Name) {
					result := a.executeTool(content.ID, content.Name, content.Input)
					toolResults = append(toolResults, result)
				} else if a.isMCPTool(content.Name) {
					result := a.executeMCPTool(content.ID, content.Name, content.Input)
					toolResults = append(toolResults, result)
				} else {
					toolResults = append(toolResults, anthropic.NewToolResultBlock(content.ID, "tool not found", true))
				}
			}
		}
		if len(toolResults) == 0 {
			readUserInput = true
			continue
		}
		readUserInput = false
		conversation = append(conversation, anthropic.NewUserMessage(toolResults...))
	}

	return nil
}

func (a *Agent) runInference(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	anthropicTools := []anthropic.ToolUnionParam{}
	for _, tool := range a.tools {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: tool.InputSchema,
			},
		})
	}
	for _, tool := range a.mcpManager.Tools() {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: anthropic.ToolInputSchemaParam{
					Type:       "object",
					Properties: tool.InputSchema.Properties,
				},
			},
		})
	}

	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: int64(1024),
		Messages:  conversation,
		Tools:     anthropicTools,
	})
	return message, err
}

func (a *Agent) isLocalTool(name string) bool {
	for _, tool := range a.tools {
		if tool.Name == name {
			return true
		}
	}
	return false
}

func (a *Agent) isMCPTool(name string) bool {
	for _, tool := range a.mcpManager.Tools() {
		if tool.Name == name {
			return true
		}
	}
	return false
}

func (a *Agent) executeTool(id, name string, input json.RawMessage) anthropic.ContentBlockParamUnion {
	var toolDef tool.ToolDefinition
	for _, tool := range a.tools {
		if tool.Name == name {
			toolDef = tool
			break
		}
	}

	fmt.Printf("\u001b[92mtool\u001b[0m: %s(%s)\n", name, input)
	response, err := toolDef.Function(input)
	if err != nil {
		return anthropic.NewToolResultBlock(id, err.Error(), true)
	}
	return anthropic.NewToolResultBlock(id, response, false)
}

func (a *Agent) executeMCPTool(id, name string, input json.RawMessage) anthropic.ContentBlockParamUnion {
	fmt.Printf("\u001b[92mtool\u001b[0m: %s(%s)\n", name, input)
	response, err := a.mcpManager.InvokeTool(name, input)
	if err != nil {
		return anthropic.NewToolResultBlock(id, err.Error(), true)
	}
	return anthropic.NewToolResultBlock(id, response, false)
}
