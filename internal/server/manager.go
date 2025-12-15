package server

import (
	"fmt"
	"gomoco/internal/models"
	"gomoco/internal/storage"
	"log"
	"sync"

	"github.com/google/uuid"
)

// Manager manages all mock servers
type Manager struct {
	mu      sync.RWMutex
	mocks   map[string]*models.MockAPI
	servers map[string]Server
	storage *storage.Storage
}

// Server interface for mock servers
type Server interface {
	Start() error
	Stop() error
	IsRunning() bool
}

// NewManager creates a new manager instance
func NewManager() *Manager {
	store, err := storage.NewStorage()
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	m := &Manager{
		mocks:   make(map[string]*models.MockAPI),
		servers: make(map[string]Server),
		storage: store,
	}

	// Load existing mocks from storage
	if err := m.loadFromStorage(); err != nil {
		log.Printf("Warning: Failed to load mocks from storage: %v", err)
	}

	return m
}

// Create creates a new mock API
func (m *Manager) Create(req *models.CreateMockAPIRequest) (*models.MockAPI, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate unique ID
	id := uuid.New().String()

	// Check if port is already in use
	for _, mock := range m.mocks {
		if mock.Port == req.Port && mock.Status == "running" {
			return nil, fmt.Errorf("port %d is already in use", req.Port)
		}
	}

	mock := &models.MockAPI{
		ID:                  id,
		Name:                req.Name,
		Port:                req.Port,
		Protocol:            req.Protocol,
		CertFile:            req.CertFile,
		KeyFile:             req.KeyFile,
		FTPMode:             req.FTPMode,
		FTPRootDir:          req.FTPRootDir,
		FTPUser:             req.FTPUser,
		FTPPass:             req.FTPPass,
		FTPPassivePortRange: req.FTPPassivePortRange,
		SFTPRootDir:         req.SFTPRootDir,
		SFTPUser:            req.SFTPUser,
		SFTPPass:            req.SFTPPass,
		SFTPHostKey:         req.SFTPHostKey,
		SFTPPrivateKey:      req.SFTPPrivateKey,
		Content:             req.Content,
		Charset:             req.Charset,
		Path:                req.Path,
		Method:              req.Method,
		Status:              "stopped",
	}

	m.mocks[id] = mock

	// Start the server
	if err := m.startServer(mock); err != nil {
		delete(m.mocks, id)
		return nil, err
	}

	mock.Status = "running"

	// Save to storage
	if err := m.saveToStorage(); err != nil {
		log.Printf("Warning: Failed to save mocks to storage: %v", err)
	}

	return mock, nil
}

// Get retrieves a mock API by ID
func (m *Manager) Get(id string) (*models.MockAPI, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	mock, exists := m.mocks[id]
	if !exists {
		return nil, fmt.Errorf("mock API not found")
	}

	return mock, nil
}

// List returns all mock APIs
func (m *Manager) List() []*models.MockAPI {
	m.mu.RLock()
	defer m.mu.RUnlock()

	mocks := make([]*models.MockAPI, 0, len(m.mocks))
	for _, mock := range m.mocks {
		mocks = append(mocks, mock)
	}

	return mocks
}

// Update updates a mock API
func (m *Manager) Update(id string, req *models.UpdateMockAPIRequest) (*models.MockAPI, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	mock, exists := m.mocks[id]
	if !exists {
		return nil, fmt.Errorf("mock API not found")
	}

	// Update fields
	if req.Name != "" {
		mock.Name = req.Name
	}
	if req.Content != "" {
		mock.Content = req.Content
	}
	if req.Charset != "" {
		mock.Charset = req.Charset
	}
	if req.Path != "" {
		mock.Path = req.Path
	}
	if req.Method != "" {
		mock.Method = req.Method
	}

	// Restart server if running
	if mock.Status == "running" {
		if err := m.stopServer(id); err != nil {
			return nil, err
		}
		if err := m.startServer(mock); err != nil {
			return nil, err
		}
	}

	// Save to storage
	if err := m.saveToStorage(); err != nil {
		log.Printf("Warning: Failed to save mocks to storage: %v", err)
	}

	return mock, nil
}

// Delete deletes a mock API
func (m *Manager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	mock, exists := m.mocks[id]
	if !exists {
		return fmt.Errorf("mock API not found")
	}

	// Stop server if running
	if mock.Status == "running" {
		if err := m.stopServer(id); err != nil {
			return err
		}
	}

	delete(m.mocks, id)

	// Save to storage
	if err := m.saveToStorage(); err != nil {
		log.Printf("Warning: Failed to save mocks to storage: %v", err)
	}

	return nil
}

// startServer starts a mock server
func (m *Manager) startServer(mock *models.MockAPI) error {
	var server Server
	var err error

	switch mock.Protocol {
	case models.ProtocolHTTP, models.ProtocolHTTPS:
		server, err = NewHTTPServer(mock)
	case models.ProtocolTCP:
		server, err = NewTCPServer(mock)
	case models.ProtocolFTP:
		server, err = NewFTPServer(mock)
	case models.ProtocolSFTP:
		server, err = NewSFTPServer(mock)
	default:
		return fmt.Errorf("unsupported protocol: %s", mock.Protocol)
	}

	if err != nil {
		return err
	}

	if err := server.Start(); err != nil {
		return err
	}

	m.servers[mock.ID] = server
	return nil
}

// stopServer stops a mock server
func (m *Manager) stopServer(id string) error {
	server, exists := m.servers[id]
	if !exists {
		return fmt.Errorf("server not found")
	}

	if err := server.Stop(); err != nil {
		return err
	}

	delete(m.servers, id)
	return nil
}

// loadFromStorage loads mocks from storage and starts them
func (m *Manager) loadFromStorage() error {
	mocks, err := m.storage.Load()
	if err != nil {
		return err
	}

	for _, mock := range mocks {
		mock.Status = "stopped"
		m.mocks[mock.ID] = mock

		// Try to start the server
		if err := m.startServer(mock); err != nil {
			log.Printf("Warning: Failed to start mock %s (%s): %v", mock.Name, mock.ID, err)
			continue
		}
		mock.Status = "running"
		log.Printf("Loaded and started mock: %s (port %d)", mock.Name, mock.Port)
	}

	return nil
}

// saveToStorage saves all mocks to storage
func (m *Manager) saveToStorage() error {
	mocks := make([]*models.MockAPI, 0, len(m.mocks))
	for _, mock := range m.mocks {
		mocks = append(mocks, mock)
	}
	return m.storage.Save(mocks)
}
