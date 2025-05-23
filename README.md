# agent-go

A Go-based agent inspired by [How to Build an Agent](https://ampcode.com/how-to-build-an-agent). Integrates the Model Context Protocol (MCP) with autonomous agents.

## Demo

![Demo of agent-go](assets/demo-speedup.gif)

## Project Structure

- `pkg/agent/` — Agent core logic and orchestration
- `pkg/mcp/` — Run and manage MCP servers
- `pkg/tool/` — Tool implementations for agent capabilities, including file operations
- `main.go` — Application entry point

## Getting Started

Pre-requisites:

- Anthropic API key
- Docker running
- Update the `config.json`

```bash
go mod download
export ANTHROPIC_API_KEY=""
go run main.go
```

<!--
LINKS:
- https://github.com/modelcontextprotocol/servers
- https://github.com/mark3labs/mcp-go
- https://github.com/metoro-io/mcp-golang
- https://github.com/llmcontext/gomcp
-->

<!--
TODO:
- Try playing with more MCP servers
- Remove dependence on Anthropic API?
  - Primarily used for `client.Messages.New()`. Accepts an existing conversation and returns a response message.
  - https://pkg.go.dev/github.com/ollama/ollama/api
- Try integrating with local Ollama models? Do local Ollama models implement the Anthropic API?
-->

<!--
DONE:
- asciinema recording
- Running and managing MCP Servers. Listing Tools, Registering Tools, Running tools
- Cursor + MCP
- Finish tutorial
-->