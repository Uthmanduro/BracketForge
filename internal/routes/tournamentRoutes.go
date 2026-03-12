package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/config"
	"github.com/uthmanduro/BracketForge/internal/handler"
	"github.com/uthmanduro/BracketForge/internal/middleware"
)

func RegisterTournamentRoutes(r *gin.RouterGroup, cfg *config.Config, handler *handler.TournamentHandler) {
	tournamentGroup := r.Group("/tournaments")
	tournamentGroup.Use(middleware.AuthMiddleware(cfg))
	{
		tournamentGroup.POST("/", handler.CreateTournament)
		tournamentGroup.GET("/:id", handler.GetTournamentByID)
		tournamentGroup.PUT("/:id", handler.UpdateTournament)
		tournamentGroup.DELETE("/:id", handler.DeleteTournament)
	}
}