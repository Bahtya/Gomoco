package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"gomoco/internal/models"
	"gomoco/internal/utils"
)

// TCPServer represents a TCP mock server
type TCPServer struct {
	mock     *models.MockAPI
	listener net.Listener
	wg       sync.WaitGroup
	stopChan chan struct{}
}

// NewTCPServer creates a new TCP server
func NewTCPServer(mock *models.MockAPI) (*TCPServer, error) {
	return &TCPServer{
		mock:     mock,
		stopChan: make(chan struct{}),
	}, nil
}

// Start starts the TCP server
func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.mock.Port))
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %v", err)
	}

	s.listener = listener

	s.wg.Add(1)
	go s.acceptConnections()

	return nil
}

// acceptConnections accepts and handles incoming connections
func (s *TCPServer) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.stopChan:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.stopChan:
					return
				default:
					fmt.Printf("TCP accept error on port %d: %v\n", s.mock.Port, err)
					continue
				}
			}

			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

// handleConnection handles a single TCP connection
func (s *TCPServer) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	// Read incoming data (optional, for logging)
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		fmt.Printf("TCP read error: %v\n", err)
	}

	// Convert content to appropriate charset
	content, err := utils.ConvertCharset(s.mock.Content, s.mock.Charset)
	if err != nil {
		fmt.Printf("Charset conversion error: %v\n", err)
		return
	}

	// Send response
	_, err = conn.Write(content)
	if err != nil {
		fmt.Printf("TCP write error: %v\n", err)
	}
}

// Stop stops the TCP server
func (s *TCPServer) Stop() error {
	if s.listener == nil {
		return nil
	}

	close(s.stopChan)
	s.listener.Close()
	s.wg.Wait()

	return nil
}

// IsRunning checks if the server is running
func (s *TCPServer) IsRunning() bool {
	return s.listener != nil
}
