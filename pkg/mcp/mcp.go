package mcp

import (
	"context"
	"fmt"
	"sync"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type Server struct {
	Name   string
	config ServerConfig

	client *client.Client
	tools  []mcp.Tool

	mu sync.Mutex
}

func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	var env []string
	for k, v := range s.config.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	s.client, err = client.NewStdioMCPClient(s.config.Command, env, s.config.Args...)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{Name: "agent-go", Version: "0.0.1"}
	_, err = s.client.Initialize(context.Background(), initRequest)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	tools, err := s.client.ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		return fmt.Errorf("failed to list tools: %w", err)
	}
	s.tools = tools.Tools

	toolNames := make([]string, 0, len(s.tools))
	for _, tool := range s.tools {
		toolNames = append(toolNames, tool.Name)
	}
	fmt.Printf("Initialized MCP Server '%s' with tools: %v ...\n", s.Name, toolNames)

	return nil
}

func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client == nil {
		return nil
	}

	if err := s.client.Close(); err != nil {
		return fmt.Errorf("failed to stop: %w", err)
	}

	s.client = nil
	return nil
}

type Manager struct {
	servers map[string]*Server
	mu      sync.RWMutex
}

func NewManager(configPath string) (*Manager, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	manager := &Manager{
		servers: make(map[string]*Server),
	}

	for name, cfg := range config.MCPServers {
		manager.servers[name] = &Server{
			Name:   name,
			config: cfg,
		}
	}

	return manager, nil
}

func (m *Manager) StartAll() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, server := range m.servers {
		if err := server.Start(); err != nil {
			return fmt.Errorf("starting %s: %w", name, err)
		}
	}
	return nil
}

func (m *Manager) StopAll() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, server := range m.servers {
		if err := server.Stop(); err != nil {
			return fmt.Errorf("stopping %s: %w", name, err)
		}
	}
	return nil
}
