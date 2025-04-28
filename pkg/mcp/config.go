package mcp

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

type ServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) || len(data) == 0 {
		if os.IsNotExist(err) {
			fmt.Println("[MCP] Config file not found: starting with no servers configured.")
		} else {
			fmt.Println("[MCP] Config file is empty: starting with no servers configured.")
		}
		return &Config{MCPServers: make(map[string]ServerConfig)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &config, nil
}
