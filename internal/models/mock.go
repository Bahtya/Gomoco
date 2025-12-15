package models

// Protocol types
const (
	ProtocolHTTP  = "http"
	ProtocolHTTPS = "https"
	ProtocolTCP   = "tcp"
	ProtocolFTP   = "ftp"
	ProtocolSFTP  = "sftp"
)

// FTP mode types
const (
	FTPModeActive  = "active"
	FTPModePassive = "passive"
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
	Protocol string `json:"protocol" yaml:"protocol" binding:"required,oneof=http https tcp ftp sftp"`
	CertFile string `json:"cert_file,omitempty" yaml:"cert_file,omitempty"` // HTTPS certificate file path
	KeyFile  string `json:"key_file,omitempty" yaml:"key_file,omitempty"`   // HTTPS private key file path
	// FTP specific fields
	FTPMode             string `json:"ftp_mode,omitempty" yaml:"ftp_mode,omitempty"`                             // active or passive
	FTPRootDir          string `json:"ftp_root_dir,omitempty" yaml:"ftp_root_dir,omitempty"`                     // FTP root directory
	FTPUser             string `json:"ftp_user,omitempty" yaml:"ftp_user,omitempty"`                             // FTP username
	FTPPass             string `json:"ftp_pass,omitempty" yaml:"ftp_pass,omitempty"`                             // FTP password
	FTPPassivePortRange string `json:"ftp_passive_port_range,omitempty" yaml:"ftp_passive_port_range,omitempty"` // Passive port range (e.g., "50000-50100")
	// SFTP specific fields
	SFTPRootDir    string `json:"sftp_root_dir,omitempty" yaml:"sftp_root_dir,omitempty"`       // SFTP root directory
	SFTPUser       string `json:"sftp_user,omitempty" yaml:"sftp_user,omitempty"`               // SFTP username
	SFTPPass       string `json:"sftp_pass,omitempty" yaml:"sftp_pass,omitempty"`               // SFTP password
	SFTPHostKey    string `json:"sftp_host_key,omitempty" yaml:"sftp_host_key,omitempty"`       // SFTP host key file path
	SFTPPrivateKey string `json:"sftp_private_key,omitempty" yaml:"sftp_private_key,omitempty"` // SFTP private key file path (optional)
	Content        string `json:"content" yaml:"content"`
	Charset        string `json:"charset" yaml:"charset" binding:"required,oneof=UTF-8 GBK"`
	Path           string `json:"path,omitempty" yaml:"path,omitempty"`     // Only for HTTP protocol
	Method         string `json:"method,omitempty" yaml:"method,omitempty"` // Only for HTTP protocol (GET, POST, etc.)
	Status         string `json:"status" yaml:"-"`                          // running, stopped
}

// CreateMockAPIRequest represents the request to create a mock API
type CreateMockAPIRequest struct {
	Name     string `json:"name" binding:"required"`
	Port     int    `json:"port" binding:"required,min=1,max=65535"`
	Protocol string `json:"protocol" binding:"required,oneof=http https tcp ftp sftp"`
	CertFile string `json:"cert_file,omitempty"` // HTTPS certificate file path
	KeyFile  string `json:"key_file,omitempty"`  // HTTPS private key file path
	// FTP specific fields
	FTPMode             string `json:"ftp_mode,omitempty"`
	FTPRootDir          string `json:"ftp_root_dir,omitempty"`
	FTPUser             string `json:"ftp_user,omitempty"`
	FTPPass             string `json:"ftp_pass,omitempty"`
	FTPPassivePortRange string `json:"ftp_passive_port_range,omitempty"`
	// SFTP specific fields
	SFTPRootDir    string `json:"sftp_root_dir,omitempty"`
	SFTPUser       string `json:"sftp_user,omitempty"`
	SFTPPass       string `json:"sftp_pass,omitempty"`
	SFTPHostKey    string `json:"sftp_host_key,omitempty"`
	SFTPPrivateKey string `json:"sftp_private_key,omitempty"`
	Content        string `json:"content"`
	Charset        string `json:"charset" binding:"required,oneof=UTF-8 GBK"`
	Path           string `json:"path,omitempty"`
	Method         string `json:"method,omitempty"`
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
	// FTP specific fields
	FTPMode             string `json:"ftp_mode,omitempty"`
	FTPRootDir          string `json:"ftp_root_dir,omitempty"`
	FTPUser             string `json:"ftp_user,omitempty"`
	FTPPass             string `json:"ftp_pass,omitempty"`
	FTPPassivePortRange string `json:"ftp_passive_port_range,omitempty"`
	// SFTP specific fields
	SFTPRootDir    string `json:"sftp_root_dir,omitempty"`
	SFTPUser       string `json:"sftp_user,omitempty"`
	SFTPPass       string `json:"sftp_pass,omitempty"`
	SFTPHostKey    string `json:"sftp_host_key,omitempty"`
	SFTPPrivateKey string `json:"sftp_private_key,omitempty"`
}
