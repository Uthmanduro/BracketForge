package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/service"
)

type PlayerHandler struct {
	PlayerService *service.PlayerService
}

type CreatePlayerRequest struct {
	Name         string `json:"name" binding:"required"`
	TournamentID uint   `json:"tournament_id" binding:"required"`
	Seed         int    `json:"seed"`
	Ranking      int    `json:"ranking"`
}

func NewPlayerHandler(playerService *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		PlayerService: playerService,
	}
}

func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var req CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	player, err := h.PlayerService.CreatePlayer(c.Request.Context(), req.Name, req.TournamentID, req.Seed, req.Ranking)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create player"})
		return
	}
	c.JSON(201, gin.H{"message": "Player created successfully", "player_id": player.ID})
}

func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid player ID"})
		return
	}

	player, err := h.PlayerService.GetPlayerByID(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "Player not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to retrieve player"})
		}
		return
	}
	c.JSON(200, player)
}

func (h *PlayerHandler) GetPlayersByTournamentID(c *gin.Context) {
	tournamentIDParam := c.Param("tournament_id")
	tournamentID, err := strconv.Atoi(tournamentIDParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid tournament ID"})
		return
	}

	players, err := h.PlayerService.GetPlayersByTournamentID(c.Request.Context(), uint(tournamentID))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve players"})
		return
	}
	c.JSON(200, players)
}

func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid player ID"})
		return
	}

	var req CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	player, err := h.PlayerService.GetPlayerByID(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "Player not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to retrieve player"})
		}
		return
	}

	player.Name = req.Name
	player.TournamentID = req.TournamentID
	player.Seed = req.Seed
	player.Ranking = &req.Ranking

	err = h.PlayerService.UpdatePlayer(c.Request.Context(), player)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update player"})
		return
	}
	c.JSON(200, gin.H{"message": "Player updated successfully"})
}

func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid player ID"})
		return
	}

	err = h.PlayerService.DeletePlayer(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "Player not found"})
		} else {
			c.JSON(500, gin.H{"error": "Failed to delete player"})
		}
		return
	}
	c.JSON(200, gin.H{"message": "Player deleted successfully"})
}



