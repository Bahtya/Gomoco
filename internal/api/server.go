package api

import (
	"embed"
	"gomoco/internal/models"
	"gomoco/internal/server"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the API server
type Server struct {
	manager     *server.Manager
	router      *gin.Engine
	staticFiles embed.FS
}

// NewServer creates a new API server
func NewServer(manager *server.Manager, staticFiles embed.FS) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	router.Use(cors.New(config))

	s := &Server{
		manager:     manager,
		router:      router,
		staticFiles: staticFiles,
	}

	s.setupRoutes()
	return s
}

// setupRoutes sets up API routes
func (s *Server) setupRoutes() {
	api := s.router.Group("/api")
	{
		api.POST("/mocks", s.createMock)
		api.GET("/mocks", s.listMocks)
		api.GET("/mocks/:id", s.getMock)
		api.PUT("/mocks/:id", s.updateMock)
		api.DELETE("/mocks/:id", s.deleteMock)

		// FTP file management
		api.GET("/mocks/:id/files", s.listFiles)
		api.GET("/mocks/:id/files/*filepath", s.downloadFile)
		api.POST("/mocks/:id/files", s.uploadFile)
		api.DELETE("/mocks/:id/files/*filepath", s.deleteFile)
	}

	// Serve embedded static files for frontend
	distFS, err := fs.Sub(s.staticFiles, "web/dist")
	if err != nil {
		panic(err)
	}

	// Serve assets directory
	assetsFS, err := fs.Sub(distFS, "assets")
	if err != nil {
		panic(err)
	}
	s.router.StaticFS("/assets", http.FS(assetsFS))

	// Serve index.html for root and 404
	s.router.GET("/", func(c *gin.Context) {
		data, err := s.staticFiles.ReadFile("web/dist/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load index.html")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	s.router.NoRoute(func(c *gin.Context) {
		data, err := s.staticFiles.ReadFile("web/dist/index.html")
		if err != nil {
			c.String(http.StatusNotFound, "Not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})
}

// createMock creates a new mock API
func (s *Server) createMock(c *gin.Context) {
	var req models.CreateMockAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mock, err := s.manager.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mock)
}

// listMocks lists all mock APIs
func (s *Server) listMocks(c *gin.Context) {
	mocks := s.manager.List()
	c.JSON(http.StatusOK, mocks)
}

// getMock gets a mock API by ID
func (s *Server) getMock(c *gin.Context) {
	id := c.Param("id")
	mock, err := s.manager.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mock)
}

// updateMock updates a mock API
func (s *Server) updateMock(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateMockAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mock, err := s.manager.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mock)
}

// deleteMock deletes a mock API
func (s *Server) deleteMock(c *gin.Context) {
	id := c.Param("id")
	if err := s.manager.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mock API deleted successfully"})
}

// Run starts the API server
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
