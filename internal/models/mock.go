package models

// Protocol types
const (
	ProtocolHTTP  = "http"
	ProtocolHTTPS = "https"
	ProtocolTCP   = "tcp"
)

// Charset types
const (
	CharsetUTF8 = "UTF-8"
	CharsetGBK  = "GBK"
)

// MockAPI represents a mock API configuration
type MockAPI struct {
	ID       string `json:"id" yaml:"id" binding:"required"`
	Name     string `json:"name" yaml:"name" binding:"required"`
	Port     int    `json:"port" yaml:"port" binding:"required,min=1,max=65535"`
	Protocol string `json:"protocol" yaml:"protocol" binding:"required,oneof=http https tcp"`
	CertFile string `json:"cert_file,omitempty" yaml:"cert_file,omitempty"` // HTTPS certificate file path
	KeyFile  string `json:"key_file,omitempty" yaml:"key_file,omitempty"`   // HTTPS private key file path
	Content  string `json:"content" yaml:"content" binding:"required"`
	Charset  string `json:"charset" yaml:"charset" binding:"required,oneof=UTF-8 GBK"`
	Path     string `json:"path,omitempty" yaml:"path,omitempty"`     // Only for HTTP protocol
	Method   string `json:"method,omitempty" yaml:"method,omitempty"` // Only for HTTP protocol (GET, POST, etc.)
	Status   string `json:"status" yaml:"-"`                          // running, stopped
}

// CreateMockAPIRequest represents the request to create a mock API
type CreateMockAPIRequest struct {
	Name     string `json:"name" binding:"required"`
	Port     int    `json:"port" binding:"required,min=1,max=65535"`
	Protocol string `json:"protocol" binding:"required,oneof=http https tcp"`
	CertFile string `json:"cert_file,omitempty"` // HTTPS certificate file path
	KeyFile  string `json:"key_file,omitempty"`  // HTTPS private key file path
	Content  string `json:"content" binding:"required"`
	Charset  string `json:"charset" binding:"required,oneof=UTF-8 GBK"`
	Path     string `json:"path,omitempty"`
	Method   string `json:"method,omitempty"`
}

// UpdateMockAPIRequest represents the request to update a mock API
type UpdateMockAPIRequest struct {
	Name     string `json:"name,omitempty"`
	Content  string `json:"content,omitempty"`
	Charset  string `json:"charset,omitempty"`
	Path     string `json:"path,omitempty"`
	Method   string `json:"method,omitempty"`
	CertFile string `json:"cert_file,omitempty"`
	KeyFile  string `json:"key_file,omitempty"`
}
