package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/middleware"
	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/service"
)

type TournamentHandler struct {
	TournamentService *service.TournamentService
}

func NewTournamentHandler(tournamentService *service.TournamentService) *TournamentHandler {
	return &TournamentHandler{
		TournamentService: tournamentService,
	}
}

func (h *TournamentHandler) CreateTournament(c *gin.Context) {
	var req model.CreateTournamentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Error binding JSON: %v\n", err)
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}
	
	tournament, err := h.TournamentService.Create(middleware.OrgID(c), &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Tournament created successfully", "tournament_id": tournament.ID})
}

func (h *TournamentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	tournament, err := h.TournamentService.GetByID(id, middleware.OrgID(c))
	if err != nil {
		c.JSON(404, gin.H{"error": "Tournament not found"})
		return
	}
	c.JSON(200, tournament)
}

func (h *TournamentHandler) List(c *gin.Context) {
	list, err := h.TournamentService.List(middleware.OrgID(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

func (h *TournamentHandler) UpdateStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status" binding:"required,oneof=registration active completed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.TournamentService.TransitionStatus(c.Param("id"), middleware.OrgID(c), req.Status); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": req.Status})
}

func (h *TournamentHandler) RegisterPlayer(c *gin.Context) {
	var req model.RegisterPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	reg, err := h.TournamentService.RegisterPlayer(c.Param("id"), middleware.OrgID(c), &req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, reg)
}

func (h *TournamentHandler) ListRegistrations(c *gin.Context) {
	list, err := h.TournamentService.ListRegistrations(c.Param("id"), middleware.OrgID(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}
 
func (h *TournamentHandler) UpdateRegistrationStatus(c *gin.Context) {
	var req model.UpdateRegistrationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.TournamentService.UpdateRegistrationStatus(c.Param("id"), middleware.OrgID(c), c.Param("rid"), req.Status); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": req.Status})
}

func (h *TournamentHandler) CreateStage(c *gin.Context) {
	var req model.CreateStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	stage, err := h.TournamentService.CreateStage(c.Param("id"), middleware.OrgID(c), &req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, stage)
}
 
func (h *TournamentHandler) ListStages(c *gin.Context) {
	list, err := h.TournamentService.ListStages(c.Param("id"), middleware.OrgID(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

func (h *TournamentHandler) RunGroupDraw(c *gin.Context) {
	var req model.DrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	groups, err := h.TournamentService.RunGroupDraw(c.Param("id"), middleware.OrgID(c), c.Param("stageId"), &req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, groups)
}
 
func (h *TournamentHandler) GenerateBracket(c *gin.Context) {
	matches, err := h.TournamentService.GenerateBracket(c.Param("id"), middleware.OrgID(c), c.Param("stageId"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, matches)
}
 
func (h *TournamentHandler) ListStageMatches(c *gin.Context) {
	matches, err := h.TournamentService.ListStageMatches(c.Param("stageId"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, matches)
}
 
func (h *TournamentHandler) GetMatch(c *gin.Context) {
	detail, err := h.TournamentService.GetMatchDetail(c.Param("matchId"))
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, detail)
}
 
func (h *TournamentHandler) SubmitResult(c *gin.Context) {
	var req model.SubmitResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	detail, err := h.TournamentService.SubmitResult(c.Request.Context(), c.Param("id"), middleware.OrgID(c), c.Param("matchId"), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, detail)
}
 
func (h *TournamentHandler) CorrectResult(c *gin.Context) {
	var req model.SubmitResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	detail, err := h.TournamentService.CorrectResult(c.Request.Context(), c.Param("id"), middleware.OrgID(c), c.Param("matchId"), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, detail)
}
 
func (h *TournamentHandler) ScheduleMatch(c *gin.Context) {
	var req model.ScheduleMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.TournamentService.ScheduleMatch(c.Param("matchId"), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"scheduled_at": req.ScheduledAt})
}
 
func (h *TournamentHandler) GetStageStandings(c *gin.Context) {
	standings, err := h.TournamentService.GetStageStandings(c.Param("stageId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, standings)
}
 
func (h *TournamentHandler) GetGroupStandings(c *gin.Context) {
	standings, err := h.TournamentService.GetGroupStandings(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, standings)
}
 