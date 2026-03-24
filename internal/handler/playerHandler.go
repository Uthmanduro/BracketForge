package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/middleware"
	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/service"
)

type PlayerHandler struct {
	PlayerService *service.PlayerService
}

func NewPlayerHandler(playerService *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		PlayerService: playerService,
	}
}

func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var req model.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	orgID := middleware.OrgID(c)

	player, err := h.PlayerService.Create(orgID, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create player"})
		return
	}
	c.JSON(201, gin.H{"message": "Player created successfully", "player_id": player.ID})
}

func (h *PlayerHandler) List(c *gin.Context) {
	players, err := h.PlayerService.List(middleware.OrgID(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, players)
}

func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	id := c.Param("id")

	orgID := middleware.OrgID(c)
	player, err := h.PlayerService.GetByID(id, orgID)
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

func (h *PlayerHandler) GetPlayersByOrgID(c *gin.Context) {
	orgID := middleware.OrgID(c)
	players, err := h.PlayerService.List(orgID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve players"})
		return
	}
	c.JSON(200, players)
}

func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	id := c.Param("id")
	orgId := middleware.OrgID(c)

	var req model.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	player, err := h.PlayerService.Update(id, orgId, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update player"})
		return
	}
	c.JSON(200, gin.H{"message": "Player updated successfully", "player": player})
}

func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	id := c.Param("id")
	orgId := middleware.OrgID(c)

	err := h.PlayerService.Delete(id, orgId)
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



