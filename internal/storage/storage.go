package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gomoco/internal/models"

	"gopkg.in/yaml.v3"
)

const (
	configDir  = "config"
	configFile = "mocks.yaml"
)

// Storage handles persistence of mock APIs
type Storage struct {
	mu       sync.RWMutex
	filePath string
}

// MocksConfig represents the YAML configuration structure
type MocksConfig struct {
	Mocks []*models.MockAPI `yaml:"mocks"`
}

// NewStorage creates a new storage instance
func NewStorage() (*Storage, error) {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	filePath := filepath.Join(configDir, configFile)

	return &Storage{
		filePath: filePath,
	}, nil
}

// Load loads mock APIs from YAML file
func (s *Storage) Load() ([]*models.MockAPI, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if file exists
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return []*models.MockAPI{}, nil
	}

	// Read file
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse YAML
	var config MocksConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return config.Mocks, nil
}

// Save saves mock APIs to YAML file
func (s *Storage) Save(mocks []*models.MockAPI) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	config := MocksConfig{
		Mocks: mocks,
	}

	// Marshal to YAML
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
