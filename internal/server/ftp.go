package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gomoco/internal/models"

	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
)

// FTPServer represents an FTP server
type FTPServer struct {
	mock    *models.MockAPI
	server  *server.Server
	running bool
}

// NewFTPServer creates a new FTP server
func NewFTPServer(mock *models.MockAPI) (*FTPServer, error) {
	// Set default values
	if mock.FTPMode == "" {
		mock.FTPMode = models.FTPModePassive
	}
	if mock.FTPRootDir == "" {
		mock.FTPRootDir = filepath.Join("ftp_data", fmt.Sprintf("port_%d", mock.Port))
	}
	if mock.FTPUser == "" {
		mock.FTPUser = "admin"
	}
	if mock.FTPPass == "" {
		mock.FTPPass = "admin"
	}

	// Create FTP root directory if it doesn't exist
	if err := os.MkdirAll(mock.FTPRootDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create FTP root directory: %v", err)
	}

	// Create file driver factory
	factory := &filedriver.FileDriverFactory{
		RootPath: mock.FTPRootDir,
		Perm:     server.NewSimplePerm("user", "group"),
	}

	// Configure FTP server options
	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     mock.Port,
		Hostname: "0.0.0.0",
		Auth:     &ftpAuth{user: mock.FTPUser, pass: mock.FTPPass},
	}

	// Configure passive mode
	if mock.FTPMode == models.FTPModePassive {
		if mock.FTPPassivePortRange != "" {
			// Parse port range (e.g., "50000-50100")
			parts := strings.Split(mock.FTPPassivePortRange, "-")
			if len(parts) == 2 {
				minPort, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
				maxPort, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
				if err1 == nil && err2 == nil && minPort > 0 && maxPort > minPort {
					opts.PassivePorts = fmt.Sprintf("%d-%d", minPort, maxPort)
				}
			}
		}
		// Default passive port range if not specified
		if opts.PassivePorts == "" {
			opts.PassivePorts = "50000-50100"
		}
	}

	ftpServer := server.NewServer(opts)

	return &FTPServer{
		mock:   mock,
		server: ftpServer,
	}, nil
}

// Start starts the FTP server
func (s *FTPServer) Start() error {
	s.running = true
	go func() {
		log.Printf("Starting FTP server on port %d (mode: %s, root: %s)",
			s.mock.Port, s.mock.FTPMode, s.mock.FTPRootDir)
		if err := s.server.ListenAndServe(); err != nil {
			log.Printf("FTP server error on port %d: %v", s.mock.Port, err)
			s.running = false
		}
	}()
	return nil
}

// Stop stops the FTP server
func (s *FTPServer) Stop() error {
	if s.server != nil {
		log.Printf("Stopping FTP server on port %d", s.mock.Port)
		s.running = false
		return s.server.Shutdown()
	}
	return nil
}

// IsRunning returns whether the FTP server is running
func (s *FTPServer) IsRunning() bool {
	return s.running
}

// ftpAuth implements simple authentication
type ftpAuth struct {
	user string
	pass string
}

func (a *ftpAuth) CheckPasswd(username, password string) (bool, error) {
	return username == a.user && password == a.pass, nil
}
