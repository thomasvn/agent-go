package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

const configPath = "config.json"

type ServerConfig struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type Config struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

type Server struct {
	Name   string
	cmd    *exec.Cmd
	config ServerConfig
	mu     sync.Mutex
}

func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cmd != nil {
		return fmt.Errorf("already running")
	}

	s.cmd = exec.Command(s.config.Command, s.config.Args...)
	s.cmd.Stdout = os.Stdout
	s.cmd.Stderr = os.Stderr

	if err := s.cmd.Start(); err != nil {
		s.cmd = nil
		return fmt.Errorf("failed to start: %w", err)
	}

	return nil
}

func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cmd == nil {
		return nil
	}

	if err := s.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to stop: %w", err)
	}

	s.cmd = nil
	return nil
}

type Manager struct {
	servers map[string]*Server
	mu      sync.RWMutex
}

func NewManager() (*Manager, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
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
