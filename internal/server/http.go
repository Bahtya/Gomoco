package server

import (
	"context"
	"fmt"
	"gomoco/internal/models"
	"gomoco/internal/utils"
	"net/http"
	"time"
)

// HTTPServer represents an HTTP mock server
type HTTPServer struct {
	mock   *models.MockAPI
	server *http.Server
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(mock *models.MockAPI) (*HTTPServer, error) {
	return &HTTPServer{
		mock: mock,
	}, nil
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()

	path := s.mock.Path
	if path == "" {
		path = "/"
	}

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		// Check method if specified
		if s.mock.Method != "" && r.Method != s.mock.Method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Convert content to appropriate charset
		content, err := utils.ConvertCharset(s.mock.Content, s.mock.Charset)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Set content type based on charset
		contentType := "text/plain"
		if s.mock.Charset == models.CharsetGBK {
			contentType += "; charset=GBK"
		} else {
			contentType += "; charset=UTF-8"
		}
		w.Header().Set("Content-Type", contentType)

		w.WriteHeader(http.StatusOK)
		w.Write(content)
	})

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.mock.Port),
		Handler: mux,
	}

	go func() {
		var err error
		if s.mock.Protocol == models.ProtocolHTTPS {
			// HTTPS server
			if s.mock.CertFile == "" || s.mock.KeyFile == "" {
				fmt.Printf("HTTPS server error on port %d: certificate or key file not specified\n", s.mock.Port)
				return
			}
			err = s.server.ListenAndServeTLS(s.mock.CertFile, s.mock.KeyFile)
		} else {
			// HTTP server
			err = s.server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("%s server error on port %d: %v\n", s.mock.Protocol, s.mock.Port, err)
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)
	return nil
}

// Stop stops the HTTP server
func (s *HTTPServer) Stop() error {
	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// IsRunning checks if the server is running
func (s *HTTPServer) IsRunning() bool {
	return s.server != nil
}
