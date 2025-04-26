# agent-go

A Go-based agent implementation inspired by [How to Build an Agent](https://ampcode.com/how-to-build-an-agent). This project explores integrating Model Context Protocol (MCP) with autonomous agents.

## Project Structure

- `/pkg/mcp` - MCP (Model Context Protocol) related implementations
- `/pkg/tool` - Tool implementations for agent capabilities including file operations
- `main.go` - Main application entry point

## Getting Started

```bash
go mod download
go run main.go
```

<!--
TODO:
- Try integrating mcp servers with the agent?
  - Our agent has to be able to independently run x number of servers?
    - https://github.com/modelcontextprotocol/servers
  - https://github.com/mark3labs/mcp-go
  - https://github.com/metoro-io/mcp-golang
  - https://github.com/llmcontext/gomcp
- Build a server to run and host MCP servers? Immediately get context?
- Try integrating with local Ollama models? Do local Ollama models implement the Anthropic API?
- Play with mcp in cursor? Podman + MCP?
-->

<!--
DONE:
- Finish tutorial
-->
