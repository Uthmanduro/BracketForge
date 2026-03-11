package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/config"
	"github.com/uthmanduro/BracketForge/internal/handler"
	"github.com/uthmanduro/BracketForge/internal/routes"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
	db *gorm.DB
	orgHandler *handler.OrganizationHandler
	userHandler *handler.UserHandler
	tournamentHandler *handler.TournamentHandler
}

func NewServer(
	config *config.Config, 
	db *gorm.DB, 
	orgHandler *handler.OrganizationHandler, 
	userHandler *handler.UserHandler, 
	tournamentHandler *handler.TournamentHandler) *Server {
	return &Server{
		config: config,
		db: db,
		orgHandler: orgHandler,
		userHandler: userHandler,
		tournamentHandler: tournamentHandler,
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

	// Register organization routes
	routes.RegisterOrgRoutes(r.Group("/api"), s.orgHandler)
	routes.RegisterUserRoutes(r.Group("/api"), s.userHandler, s.config)
	routes.RegisterTournamentRoutes(r.Group("/api"), s.tournamentHandler) // Pass the tournament handler to the route registration

	return r
}