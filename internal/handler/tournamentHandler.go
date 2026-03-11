package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/service"
)

type TournamentHandler struct {
	TournamentService *service.TournamentService
}

type CreateTournamentRequest struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	Format         string `json:"format" binding:"required,oneof=single_elimination round_robin group_knockout"`
	OrganizationID uint   `json:"organization_id" binding:"required"`
	StartDate      time.Time `json:"start_date" binding:"required"`
	EndDate        time.Time `json:"end_date" binding:"required"`   
}

func NewTournamentHandler(tournamentService *service.TournamentService) *TournamentHandler {
	return &TournamentHandler{
		TournamentService: tournamentService,
	}
}

func (h *TournamentHandler) CreateTournament(c *gin.Context) {
	var req CreateTournamentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Error binding JSON: %v\n", err)
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}
	
	tournament, err := h.TournamentService.CreateTournament(
		c.Request.Context(),
		req.Name, 
		req.Description, 
		req.Format, 
		req.OrganizationID, 
		&req.StartDate, 
		&req.EndDate,
	)
	if err != nil {
		if err.Error() == "tournament with the same name already exists in the organization" {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(500, gin.H{"error": "Failed to create tournament"})
		}
		return
	}
	c.JSON(201, gin.H{"message": "Tournament created successfully", "tournament_id": tournament.ID})
}

func (h *TournamentHandler) GetTournamentByID(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid tournament ID"})
		return
	}

	tournament, err := h.TournamentService.GetTournamentByID(c.Request.Context(), uint(idUint))
	if err != nil {
		c.JSON(404, gin.H{"error": "Tournament not found"})
		return
	}
	c.JSON(200, tournament)
}

func (h *TournamentHandler) UpdateTournament(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid tournament ID"})
		return
	}

	var req CreateTournamentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	tournament, err := h.TournamentService.GetTournamentByID(c.Request.Context(), uint(idUint))
	if err != nil {
		c.JSON(404, gin.H{"error": "Tournament not found"})
		return
	}

	tournament.Name = req.Name
	tournament.Description = req.Description
	tournament.Format = model.TournamentFormat(req.Format)
	tournament.OrganizationID = req.OrganizationID
	tournament.StartDate = &req.StartDate
	tournament.EndDate = &req.EndDate

	err = h.TournamentService.UpdateTournament(c.Request.Context(), tournament)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update tournament"})
		return
	}
	c.JSON(200, gin.H{"message": "Tournament updated successfully"})
}

func (h *TournamentHandler) DeleteTournament(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid tournament ID"})
		return
	}

	err = h.TournamentService.DeleteTournament(c.Request.Context(), uint(idUint))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete tournament"})
		return
	}
	c.JSON(200, gin.H{"message": "Tournament deleted successfully"})
}