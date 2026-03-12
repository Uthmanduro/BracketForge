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
		players.POST("/",  playerHandler.CreatePlayer)
		players.GET("/:id", playerHandler.GetPlayerByID)
		players.GET("/tournaments/:tournament_id", playerHandler.GetPlayersByTournamentID)
		players.PUT("/:id", playerHandler.UpdatePlayer)
		players.DELETE("/:id", playerHandler.DeletePlayer)
	}
}