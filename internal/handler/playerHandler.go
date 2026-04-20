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

// Create godoc
// @Summary      Create player
// @Description  Creates a new player profile within the organisation. Admin or organizer only.
// @Tags         players
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body model.CreatePlayerRequest true "Player details"
// @Success      201 {object} model.SuccessResponse{Data=object{player_id=string}, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /players [post]
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var req model.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.ErrorResponse{Error: "Invalid request data"})
		return
	}

	orgID := middleware.OrgID(c)

	checkPlayer, err := h.PlayerService.GetByEmail(*req.Email, orgID)
	if err != nil && err.Error() != "record not found" {
		c.JSON(500, model.ErrorResponse{Error: "Failed to check existing player"})
		return
	}
	if err == nil && checkPlayer != nil {
		c.JSON(400, model.ErrorResponse{Error: "Player with this email already exists"})
		return
	}

	player, err := h.PlayerService.Create(orgID, &req)
	if err != nil {
		c.JSON(500, model.ErrorResponse{Error: "Failed to create player"})
		return
	}
	c.JSON(201, model.SuccessResponse{Message: "Player created successfully", Data: gin.H{"player_id": player.ID}})
}

// List godoc
// @Summary      List players
// @Description  Returns all players in the authenticated organisation
// @Tags         players
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} model.SuccessResponse{Data=object{players=[]model.Player}, Message=string}
// @Failure      401 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /players [get]
func (h *PlayerHandler) List(c *gin.Context) {
	players, err := h.PlayerService.List(middleware.OrgID(c))
	if err != nil {
		c.JSON(500, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, model.SuccessResponse{Message: "Players retrieved successfully", Data: gin.H{"players": players}})
}

// GetByID godoc
// @Summary      Get player
// @Description  Returns a single player by ID
// @Tags         players
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Player ID"
// @Success      200 {object} model.SuccessResponse{Data=object{player=model.Player}, Message=string}
// @Failure      401 {object} model.ErrorResponse
// @Failure      404 {object} model.ErrorResponse
// @Router       /players/{id} [get]
func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	id := c.Param("id")

	orgID := middleware.OrgID(c)
	player, err := h.PlayerService.GetByID(id, orgID)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(404,model.ErrorResponse{Error: "Player not found"})
		} else {
			c.JSON(500, model.ErrorResponse{Error: "Failed to retrieve player"})
		}
		return
	}
	c.JSON(200, model.SuccessResponse{Message: "Player retrieved successfully", Data: gin.H{"player": player}})
}

func (h *PlayerHandler) GetPlayersByOrgID(c *gin.Context) {
	orgID := middleware.OrgID(c)
	players, err := h.PlayerService.List(orgID)
	if err != nil {
		c.JSON(500, model.ErrorResponse{Error: "Failed to retrieve players"})
		return
	}
	c.JSON(200, players)
}

// Update godoc
// @Summary      Update player
// @Description  Updates name, email, or metadata for a player. Admin or organizer only.
// @Tags         players
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string                   true "Player ID"
// @Param        body body model.CreatePlayerRequest true "Updated player details"
// @Success      200 {object} model.SuccessResponse{Data=object{player=model.Player}, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /players/{id} [put]
func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	id := c.Param("id")
	orgId := middleware.OrgID(c)

	var req model.UpdatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.ErrorResponse{Error: "Invalid request data"})
		return
	}

	player, err := h.PlayerService.Update(id, orgId, &req)
	if err != nil {
		c.JSON(500, model.ErrorResponse{Error: "Failed to update player"})
		return
	}
	c.JSON(200, model.SuccessResponse{Message: "Player updated successfully", Data: gin.H{"player": player}})
}

// Delete godoc
// @Summary      Delete player
// @Description  Permanently deletes a player from the organisation. Admin only.
// @Tags         players
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Player ID"
// @Success      204
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /players/{id} [delete]
func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	id := c.Param("id")
	orgId := middleware.OrgID(c)

	err := h.PlayerService.Delete(id, orgId)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, model.ErrorResponse{Error: "Player not found"})
		} else {
			c.JSON(500, model.ErrorResponse{Error: "Failed to delete player"})
		}
		return
	}
	c.JSON(204, nil)
}



