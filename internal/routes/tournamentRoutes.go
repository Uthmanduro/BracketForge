package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/handler"
)

func RegisterTournamentRoutes(r *gin.RouterGroup, handler *handler.TournamentHandler) {
	tournamentGroup := r.Group("/tournaments")
	{
		tournamentGroup.POST("/", handler.CreateTournament)
		tournamentGroup.GET("/:id", handler.GetTournamentByID)
		tournamentGroup.PUT("/:id", handler.UpdateTournament)
		tournamentGroup.DELETE("/:id", handler.DeleteTournament)
	}
}