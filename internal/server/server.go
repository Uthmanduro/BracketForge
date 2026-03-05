package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/config"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
	db *gorm.DB
}

func NewServer(config *config.Config, db *gorm.DB) *Server {
	return &Server{
		config: config,
		db: db,
	}
}

func (s *Server) Start() error {
	r := s.setupRouter()

	return r.Run(":" + s.config.Port)
}

func (s *Server) setupRouter() *gin.Engine {
	r := gin.Default()

	// Add middleware, routes, etc. here
	r.GET("/health", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server healthy",
		})
	})

	return r
}