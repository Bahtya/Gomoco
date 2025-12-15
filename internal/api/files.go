package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const maxFileSize = 100 * 1024 * 1024 // 100MB

// FileInfo represents file information
type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	ModTime string `json:"mod_time"`
	Path    string `json:"path"`
}

// listFiles lists files in FTP directory
func (s *Server) listFiles(c *gin.Context) {
	id := c.Param("id")
	mock, err := s.manager.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mock API not found"})
		return
	}

	if mock.Protocol != "ftp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not an FTP mock API"})
		return
	}

	// Get path from query parameter
	subPath := c.DefaultQuery("path", "")
	fullPath := filepath.Join(mock.FTPRootDir, subPath)

	// Security check: ensure path is within root directory
	absRoot, _ := filepath.Abs(mock.FTPRootDir)
	absPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absPath, absRoot) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		relPath := filepath.Join(subPath, entry.Name())
		files = append(files, FileInfo{
			Name:    entry.Name(),
			Size:    info.Size(),
			IsDir:   entry.IsDir(),
			ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
			Path:    filepath.ToSlash(relPath),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"files":        files,
		"current_path": subPath,
		"root_dir":     mock.FTPRootDir,
	})
}

// downloadFile downloads a file from FTP directory
func (s *Server) downloadFile(c *gin.Context) {
	id := c.Param("id")
	filePath := c.Param("filepath")
	if filePath != "" && filePath[0] == '/' {
		filePath = filePath[1:]
	}

	mock, err := s.manager.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mock API not found"})
		return
	}

	if mock.Protocol != "ftp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not an FTP mock API"})
		return
	}

	fullPath := filepath.Join(mock.FTPRootDir, filePath)

	// Security check
	absRoot, _ := filepath.Abs(mock.FTPRootDir)
	absPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absPath, absRoot) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Check if file exists and is not a directory
	info, err := os.Stat(fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	if info.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot download a directory"})
		return
	}

	c.FileAttachment(fullPath, filepath.Base(filePath))
}

// uploadFile uploads a file to FTP directory
func (s *Server) uploadFile(c *gin.Context) {
	id := c.Param("id")
	mock, err := s.manager.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mock API not found"})
		return
	}

	if mock.Protocol != "ftp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not an FTP mock API"})
		return
	}

	// Get upload path from form
	uploadPath := c.PostForm("path")

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Check file size (100MB limit)
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("File size exceeds 100MB limit (size: %.2f MB)", float64(file.Size)/(1024*1024)),
		})
		return
	}

	// Construct full path
	fullPath := filepath.Join(mock.FTPRootDir, uploadPath, file.Filename)

	// Security check
	absRoot, _ := filepath.Abs(mock.FTPRootDir)
	absPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absPath, absRoot) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	// Save file
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": file.Filename,
		"size":     file.Size,
		"path":     filepath.ToSlash(filepath.Join(uploadPath, file.Filename)),
	})
}

// deleteFile deletes a file from FTP directory
func (s *Server) deleteFile(c *gin.Context) {
	id := c.Param("id")
	filePath := c.Param("filepath")
	if filePath != "" && filePath[0] == '/' {
		filePath = filePath[1:]
	}

	mock, err := s.manager.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mock API not found"})
		return
	}

	if mock.Protocol != "ftp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not an FTP mock API"})
		return
	}

	fullPath := filepath.Join(mock.FTPRootDir, filePath)

	// Security check
	absRoot, _ := filepath.Abs(mock.FTPRootDir)
	absPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absPath, absRoot) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Check if file/directory exists
	info, err := os.Stat(fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Delete file or directory
	if info.IsDir() {
		err = os.RemoveAll(fullPath)
	} else {
		err = os.Remove(fullPath)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
		"path":    filePath,
	})
}
