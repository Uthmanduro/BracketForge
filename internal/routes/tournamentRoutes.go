package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/config"
	"github.com/uthmanduro/BracketForge/internal/handler"
	"github.com/uthmanduro/BracketForge/internal/middleware"
)

func RegisterTournamentRoutes(r *gin.RouterGroup, cfg *config.Config, handler *handler.TournamentHandler) {
	base := r.Group("/tournaments")
	base.Use(middleware.AuthMiddleware(cfg))

	adminOrOrg := middleware.RequireRole("admin", "organizer")

	// ── Tournament CRUD ────────────────────────────────────────────────
	base.POST("", adminOrOrg, handler.CreateTournament)
	base.GET("", handler.List)
	base.GET("/:id", handler.GetByID)
	base.PATCH("/:id/status", adminOrOrg, handler.UpdateStatus)

	// ── Registrations ──────────────────────────────────────────────────
	base.POST("/:id/registrations", adminOrOrg, handler.RegisterPlayer)
	base.GET("/:id/registrations", handler.ListRegistrations)
	base.PATCH("/:id/registrations/:rid", adminOrOrg, handler.UpdateRegistrationStatus)
 
	// ── Stages ─────────────────────────────────────────────────────────
	base.POST("/:id/stages", adminOrOrg, handler.CreateStage)
	base.GET("/:id/stages", handler.ListStages)
 
	// ── Draw and bracket generation ────────────────────────────────────
	base.POST("/:id/stages/:stageId/draw", adminOrOrg, handler.RunGroupDraw)
	base.POST("/:id/stages/:stageId/bracket", adminOrOrg, handler.GenerateBracket)
 
	// ── Stage views ────────────────────────────────────────────────────
	base.GET("/:id/stages/:stageId/matches", handler.ListStageMatches)
	base.GET("/:id/stages/:stageId/standings", handler.GetStageStandings)
 
	// ── Group standings ────────────────────────────────────────────────
	base.GET("/:id/groups/:groupId/standings", handler.GetGroupStandings)
 
	// ── Match operations ───────────────────────────────────────────────
	base.GET("/:id/matches/:matchId", handler.GetMatch)
	base.POST("/:id/matches/:matchId/result", adminOrOrg, handler.SubmitResult)
	base.PUT("/:id/matches/:matchId/result", middleware.RequireRole("admin"), handler.CorrectResult)
	base.PATCH("/:id/matches/:matchId/schedule", adminOrOrg, handler.ScheduleMatch)
}