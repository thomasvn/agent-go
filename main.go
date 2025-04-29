package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"agent/pkg/agent"
	"agent/pkg/mcp"
	"agent/pkg/tool"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	configPath := flag.String("config", "myconfig.json", "path to MCP config file")
	flag.Parse()

	mcpManager, err := mcp.NewManager(*configPath)
	if err != nil {
		fmt.Printf("Error initializing MCP Manager: %s\n", err)
		return
	}
	err = mcpManager.StartAll()
	if err != nil {
		fmt.Printf("Error starting MCP servers: %s\n", err)
		return
	}

	client := anthropic.NewClient()
	tools := []tool.ToolDefinition{tool.ReadFileDefinition, tool.ListFilesDefinition, tool.EditFileDefinition}

	agent := agent.NewAgent(&client, getUserMessage, tools, mcpManager)
	err = agent.Run(context.TODO())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

func getUserMessage() (string, bool) {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", false
	}
	return scanner.Text(), true
}
