package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"gomoco/internal/models"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SFTPServer represents an SFTP server
type SFTPServer struct {
	mock     *models.MockAPI
	listener net.Listener
	running  bool
	stopChan chan struct{}
}

// NewSFTPServer creates a new SFTP server
func NewSFTPServer(mock *models.MockAPI) (*SFTPServer, error) {
	// Set default values
	if mock.SFTPRootDir == "" {
		mock.SFTPRootDir = filepath.Join("sftp_data", fmt.Sprintf("port_%d", mock.Port))
	}
	if mock.SFTPUser == "" {
		mock.SFTPUser = "admin"
	}
	if mock.SFTPPass == "" {
		mock.SFTPPass = "admin"
	}

	// Create SFTP root directory if it doesn't exist
	if err := os.MkdirAll(mock.SFTPRootDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create SFTP root directory: %v", err)
	}

	// Generate or load host key
	hostKeyPath := mock.SFTPHostKey
	if hostKeyPath == "" {
		hostKeyPath = filepath.Join("sftp_keys", fmt.Sprintf("host_key_%d", mock.Port))
	}

	return &SFTPServer{
		mock:     mock,
		stopChan: make(chan struct{}),
	}, nil
}

// Start starts the SFTP server
func (s *SFTPServer) Start() error {
	// Configure SSH server
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == s.mock.SFTPUser && string(pass) == s.mock.SFTPPass {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	// Load or generate host key
	hostKey, err := s.getOrCreateHostKey()
	if err != nil {
		return fmt.Errorf("failed to get host key: %v", err)
	}
	config.AddHostKey(hostKey)

	// Start listening
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.mock.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %v", s.mock.Port, err)
	}
	s.listener = listener
	s.running = true

	log.Printf("Starting SFTP server on port %d (root: %s, user: %s)",
		s.mock.Port, s.mock.SFTPRootDir, s.mock.SFTPUser)

	// Accept connections
	go func() {
		for {
			select {
			case <-s.stopChan:
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					if s.running {
						log.Printf("SFTP server accept error on port %d: %v", s.mock.Port, err)
					}
					continue
				}

				go s.handleConnection(conn, config)
			}
		}
	}()

	return nil
}

// Stop stops the SFTP server
func (s *SFTPServer) Stop() error {
	if s.listener != nil {
		log.Printf("Stopping SFTP server on port %d", s.mock.Port)
		s.running = false
		close(s.stopChan)
		return s.listener.Close()
	}
	return nil
}

// IsRunning returns whether the SFTP server is running
func (s *SFTPServer) IsRunning() bool {
	return s.running
}

// handleConnection handles a single SSH connection
func (s *SFTPServer) handleConnection(conn net.Conn, config *ssh.ServerConfig) {
	defer conn.Close()

	// Perform SSH handshake
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Printf("SFTP SSH handshake failed: %v", err)
		return
	}
	defer sshConn.Close()

	// Discard all global requests
	go ssh.DiscardRequests(reqs)

	// Handle channels
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("SFTP channel accept failed: %v", err)
			continue
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				ok := false
				switch req.Type {
				case "subsystem":
					if string(req.Payload[4:]) == "sftp" {
						ok = true
					}
				}
				req.Reply(ok, nil)
			}
		}(requests)

		// Create SFTP server
		// Use absolute path to avoid Windows path issues
		absRootDir, err := filepath.Abs(s.mock.SFTPRootDir)
		if err != nil {
			log.Printf("SFTP failed to get absolute path: %v", err)
			channel.Close()
			continue
		}

		server, err := sftp.NewServer(
			channel,
			sftp.WithServerWorkingDirectory(absRootDir),
		)
		if err != nil {
			log.Printf("SFTP server creation failed: %v", err)
			channel.Close()
			continue
		}

		if err := server.Serve(); err != nil && err != io.EOF {
			log.Printf("SFTP server error: %v", err)
		}
		server.Close()
		channel.Close()
	}
}

// getOrCreateHostKey loads or generates a host key
func (s *SFTPServer) getOrCreateHostKey() (ssh.Signer, error) {
	hostKeyPath := s.mock.SFTPHostKey
	if hostKeyPath == "" {
		hostKeyPath = filepath.Join("sftp_keys", fmt.Sprintf("host_key_%d", s.mock.Port))
	}

	// Try to load existing key
	if _, err := os.Stat(hostKeyPath); err == nil {
		keyData, err := os.ReadFile(hostKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read host key: %v", err)
		}
		return ssh.ParsePrivateKey(keyData)
	}

	// Generate new key
	log.Printf("Generating new host key for SFTP server on port %d", s.mock.Port)

	// Use a simple approach: generate RSA key
	privateKey, err := generateRSAKey(2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %v", err)
	}

	// Save the key
	if err := os.MkdirAll(filepath.Dir(hostKeyPath), 0700); err != nil {
		return nil, fmt.Errorf("failed to create key directory: %v", err)
	}

	if err := os.WriteFile(hostKeyPath, privateKey, 0600); err != nil {
		return nil, fmt.Errorf("failed to save host key: %v", err)
	}

	return ssh.ParsePrivateKey(privateKey)
}

// generateRSAKey generates an RSA private key in PEM format
func generateRSAKey(bits int) ([]byte, error) {
	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %v", err)
	}

	// Encode private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	return pem.EncodeToMemory(privateKeyPEM), nil
}
