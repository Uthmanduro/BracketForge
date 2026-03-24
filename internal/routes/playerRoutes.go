package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/config"
	"github.com/uthmanduro/BracketForge/internal/handler"
	"github.com/uthmanduro/BracketForge/internal/middleware"
)

func RegisterPlayerRoutes(rg *gin.RouterGroup, cfg *config.Config, playerHandler *handler.PlayerHandler) {
	players := rg.Group("/players")
	players.Use(middleware.AuthMiddleware(cfg))
	{
		players.POST("/", middleware.RequireRole("admin", "organizer"), playerHandler.CreatePlayer)
		players.GET("/", playerHandler.List)
		players.GET("/:id", playerHandler.GetPlayerByID)
		// players.GET("/organizations/", playerHandler.GetPlayersByOrgID)
		players.PUT("/:id", middleware.RequireRole("admin", "organizer"), playerHandler.UpdatePlayer)
		players.DELETE("/:id", middleware.RequireRole("admin"), playerHandler.DeletePlayer)
	}
}