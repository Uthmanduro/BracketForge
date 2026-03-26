package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/config"
	"github.com/uthmanduro/BracketForge/internal/handler"
	"github.com/uthmanduro/BracketForge/internal/routes"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config *config.Config
	db *gorm.DB
	orgHandler *handler.OrganizationHandler
	userHandler *handler.UserHandler
	tournamentHandler *handler.TournamentHandler
	playerHandler *handler.PlayerHandler
}

func NewServer(
	config *config.Config, 
	db *gorm.DB, 
	orgHandler *handler.OrganizationHandler, 
	userHandler *handler.UserHandler, 
	tournamentHandler *handler.TournamentHandler, 
	playerHandler *handler.PlayerHandler) *Server {
	return &Server{
		config: config,
		db: db,
		orgHandler: orgHandler,
		userHandler: userHandler,
		tournamentHandler: tournamentHandler,
		playerHandler: playerHandler,
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

	// Swagger UI — http://localhost:8080/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register organization routes
	routes.RegisterOrgRoutes(r.Group("/api/v1"), s.config, s.orgHandler)
	routes.RegisterUserRoutes(r.Group("/api/v1"), s.userHandler, s.config)
	routes.RegisterTournamentRoutes(r.Group("/api/v1"), s.config, s.tournamentHandler) // Pass the tournament handler to the route registration
	routes.RegisterPlayerRoutes(r.Group("/api/v1"), s.config, s.playerHandler) // Pass the player handler to the route registration

	return r
}