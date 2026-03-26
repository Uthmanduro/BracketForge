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

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=registration active completed"`
}

func NewTournamentHandler(tournamentService *service.TournamentService) *TournamentHandler {
	return &TournamentHandler{
		TournamentService: tournamentService,
	}
}

// Create godoc
// @Summary      Create tournament
// @Description  Creates a new tournament. Format must be one of: single_elimination, round_robin, group_knockout. BestOf must be 3 or 5.
// @Tags         tournaments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body model.CreateTournamentRequest true "Tournament configuration"
// @Success      201 {object} model.SuccessResponse{Data=object{tournament_id=string}, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /tournaments [post]
func (h *TournamentHandler) CreateTournament(c *gin.Context) {
	var req model.CreateTournamentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Error binding JSON: %v\n", err)
		c.JSON(400, model.ErrorResponse{Error: "Invalid request data"})
		return
	}
	
	tournament, err := h.TournamentService.Create(middleware.OrgID(c), &req)
	if err != nil {
		c.JSON(500, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(201, model.SuccessResponse{Message: "Tournament created successfully", Data: gin.H{"tournament_id": tournament.ID}})
}

// GetByID godoc
// @Summary      Get tournament
// @Description  Returns a single tournament by ID
// @Tags         tournaments
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Tournament ID"
// @Success      200 {object} model.SuccessResponse{Data=object{tournament=model.Tournament}, Message=string}
// @Failure      401 {object} model.ErrorResponse
// @Failure      404 {object} model.ErrorResponse
// @Router       /tournaments/{id} [get]
func (h *TournamentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	tournament, err := h.TournamentService.GetByID(id, middleware.OrgID(c))
	if err != nil {
		c.JSON(404, model.ErrorResponse{Error: "Tournament not found"})
		return
	}
	c.JSON(200, model.SuccessResponse{Message: "Tournament found", Data: gin.H{"tournament": tournament}})
}

// List godoc
// @Summary      List tournaments
// @Description  Returns all tournaments belonging to the authenticated organisation
// @Tags         tournaments
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object}  model.SuccessResponse{Data=object{tournaments=[]model.Tournament}, Message=string}
// @Failure      401 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /tournaments [get]
func (h *TournamentHandler) List(c *gin.Context) {
	list, err := h.TournamentService.List(middleware.OrgID(c))
	if err != nil {
		c.JSON(500, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, model.SuccessResponse{Message: "Tournament retrieved successfully", Data: gin.H{"tournaments": list}})
}

// UpdateStatus godoc
// @Summary      Transition tournament status
// @Description  Moves the tournament through its lifecycle: draft → registration → active → completed. Each step must be taken in order.
// @Tags         tournaments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string                     true "Tournament ID"
// @Param        body body UpdateStatusRequest  true "New status"
// @Success      200 {object} model.SuccessResponse{Data=object{status=string}, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Router       /tournaments/{id}/status [patch]
func (h *TournamentHandler) UpdateStatus(c *gin.Context) {
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.ErrorResponse{Error: err.Error()})
		return
	}
	if err := h.TournamentService.TransitionStatus(c.Param("id"), middleware.OrgID(c), req.Status); err != nil {
		c.JSON(400, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, model.SuccessResponse{Message: "Tournament status updated successfully", Data: gin.H{"status": req.Status}})
}

// RegisterPlayer godoc
// @Summary      Register player
// @Description  Registers a player into a tournament. Tournament must be in draft or registration status.
// @Tags         registrations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string                      true "Tournament ID"
// @Param        body body model.RegisterPlayerRequest true "Player registration"
// @Success      201 {object} model.SuccessResponse{Data=object{player_registration=model.TournamentRegistration}, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Router       /tournaments/{id}/registrations [post]
func (h *TournamentHandler) RegisterPlayer(c *gin.Context) {
	var req model.RegisterPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.ErrorResponse{Error: err.Error()})
		return
	}
	reg, err := h.TournamentService.RegisterPlayer(c.Param("id"), middleware.OrgID(c), &req)
	if err != nil {
		c.JSON(400, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(201, model.SuccessResponse{Message: "Player registered successfully", Data: gin.H{"player_registration": reg}})
}

// ListRegistrations godoc
// @Summary      List registrations
// @Description  Returns all registrations for a tournament including player names
// @Tags         registrations
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Tournament ID"
// @Success      200 {object}  model.SuccessResponse{Data=object{registrations=[]model.TournamentRegistration}, Message=string}
// @Failure      401 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /tournaments/{id}/registrations [get]
func (h *TournamentHandler) ListRegistrations(c *gin.Context) {
	list, err := h.TournamentService.ListRegistrations(c.Param("id"), middleware.OrgID(c))
	if err != nil {
		c.JSON(500, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, model.SuccessResponse{Message: "Registrations retrieved successfully", Data: gin.H{"registrations": list}})
}

// UpdateRegistrationStatus godoc
// @Summary      Update registration status
// @Description  Updates a player's registration status. Setting withdrawn or disqualified while the tournament is active triggers the walkover engine for all pending matches.
// @Tags         registrations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string                                true "Tournament ID"
// @Param        rid  path string                                true "Registration ID"
// @Param        body body model.UpdateRegistrationStatusRequest true "New status"
// @Success      200 {object} model.SuccessResponse{Data=object{status=string}, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Router       /tournaments/{id}/registrations/{rid} [patch]
func (h *TournamentHandler) UpdateRegistrationStatus(c *gin.Context) {
	var req model.UpdateRegistrationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.ErrorResponse{Error: err.Error()})
		return
	}
	if err := h.TournamentService.UpdateRegistrationStatus(c.Param("id"), middleware.OrgID(c), c.Param("rid"), req.Status); err != nil {
		c.JSON(400, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, model.SuccessResponse{Message: "Registration status updated successfully", Data: gin.H{"status": req.Status}})
}

// CreateStage godoc
// @Summary      Create stage
// @Description  Creates a stage within a tournament. Type must be one of: group, knockout, round_robin. Stage order controls the sequence (1 = first played).
// @Tags         stages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string                  true "Tournament ID"
// @Param        body body model.CreateStageRequest true "Stage configuration"
// @Success      201 {object} model.SuccessResponse{Data=object{stage=model.Stage}, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Router       /tournaments/{id}/stages [post]
func (h *TournamentHandler) CreateStage(c *gin.Context) {
	var req model.CreateStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.ErrorResponse{Error: err.Error()})
		return
	}
	stage, err := h.TournamentService.CreateStage(c.Param("id"), middleware.OrgID(c), &req)
	if err != nil {
		c.JSON(400, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(201, model.SuccessResponse{Message: "Stage created successfully", Data: gin.H{"stage": stage}})
}

// ListStages godoc
// @Summary      List stages
// @Description  Returns all stages for a tournament ordered by stage_order
// @Tags         stages
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Tournament ID"
// @Success      200 {object}  model.SuccessResponse{Data=object{stages=[]model.Stage}, Message=string}
// @Failure      401 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /tournaments/{id}/stages [get]
func (h *TournamentHandler) ListStages(c *gin.Context) {
	list, err := h.TournamentService.ListStages(c.Param("id"), middleware.OrgID(c))
	if err != nil {
		c.JSON(500, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, model.SuccessResponse{Message: "Stages retrieved successfully", Data: gin.H{"stages": list}})
}

// RunGroupDraw godoc
// @Summary      Run group draw
// @Description  Distributes registered players into groups using serpentine seeding and generates all round-robin matches within each group. Only valid for stages of type group.
// @Tags         bracket
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path string             true "Tournament ID"
// @Param        stageId path string             true "Stage ID"
// @Param        body    body model.DrawRequest  true "Number of groups"
// @Success      201 {object} model.SuccessResponse{Data=object{groups=[]model.Group}, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Router       /tournaments/{id}/stages/{stageId}/draw [post]
func (h *TournamentHandler) RunGroupDraw(c *gin.Context) {
	var req model.DrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.ErrorResponse{Error: err.Error()})
		return
	}
	groups, err := h.TournamentService.RunGroupDraw(c.Param("id"), middleware.OrgID(c), c.Param("stageId"), &req)
	if err != nil {
		c.JSON(400, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(201, model.SuccessResponse{Message: "Group draw completed successfully", Data: gin.H{"groups": groups}})
}

// GenerateBracket godoc
// @Summary      Generate bracket
// @Description  Generates the match bracket for a knockout or round-robin stage. BYE matches are auto-completed. Use /draw for group stages instead.
// @Tags         bracket
// @Produce      json
// @Security     BearerAuth
// @Param        id      path string true "Tournament ID"
// @Param        stageId path string true "Stage ID"
// @Success      201 {object}  model.SuccessResponse{Data=object{matches=[]model.Match}, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Router       /tournaments/{id}/stages/{stageId}/bracket [post]
func (h *TournamentHandler) GenerateBracket(c *gin.Context) {
	matches, err := h.TournamentService.GenerateBracket(c.Param("id"), middleware.OrgID(c), c.Param("stageId"))
	if err != nil {
		c.JSON(400, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(201, model.SuccessResponse{Message: "Bracket generated successfully", Data: gin.H{"matches": matches}})
}

// ListStageMatches godoc
// @Summary      List stage matches
// @Description  Returns all matches for a stage ordered by round and position
// @Tags         matches
// @Produce      json
// @Security     BearerAuth
// @Param        id      path string true "Tournament ID"
// @Param        stageId path string true "Stage ID"
// @Success      200 {object} model.SuccessResponse{Data=object{matches=[]model.Match}, Message=string}
// @Failure      401 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /tournaments/{id}/stages/{stageId}/matches [get]
func (h *TournamentHandler) ListStageMatches(c *gin.Context) {
	matches, err := h.TournamentService.ListStageMatches(c.Param("stageId"))
	if err != nil {
		c.JSON(500, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, model.SuccessResponse{Message: "Stage matches retrieved successfully", Data: gin.H{"matches": matches}})
}

// GetMatch godoc
// @Summary      Get match detail
// @Description  Returns a match with its full participant list and set-by-set scores
// @Tags         matches
// @Produce      json
// @Security     BearerAuth
// @Param        id      path string true "Tournament ID"
// @Param        matchId path string true "Match ID"
// @Success      200 {object} model.SuccessResponse{Data=object{match=model.MatchDetail}, Message=string}
// @Failure      401 {object} model.ErrorResponse
// @Failure      404 {object} model.ErrorResponse
// @Router       /tournaments/{id}/matches/{matchId} [get]
func (h *TournamentHandler) GetMatch(c *gin.Context) {
	detail, err := h.TournamentService.GetMatchDetail(c.Param("matchId"))
	if err != nil {
		c.JSON(404, model.ErrorResponse{Error: "not found"})
		return
	}
	c.JSON(200, model.SuccessResponse{Message: "Match detail retrieved successfully", Data: gin.H{"match": detail}})
}

// SubmitResult godoc
// @Summary      Submit match result
// @Description  Records the full set-by-set result for a match. Validates set count against best_of, updates standings for group/round-robin stages, and advances the winner to the next bracket slot atomically.
// @Tags         matches
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path string                    true "Tournament ID"
// @Param        matchId path string                    true "Match ID"
// @Param        body    body model.SubmitResultRequest true "Set scores"
// @Success      200 {object} model.SuccessResponse{Data=object{result_detail=model.MatchDetail}, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Router       /tournaments/{id}/matches/{matchId}/result [post]
func (h *TournamentHandler) SubmitResult(c *gin.Context) {
	var req model.SubmitResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	detail, err := h.TournamentService.SubmitResult(c.Request.Context(), c.Param("id"), middleware.OrgID(c), c.Param("matchId"), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Result submitted successfully", Data: gin.H{"result_detail": detail}})
}

// CorrectResult godoc
// @Summary      Correct match result
// @Description  Voids the previous result, reverses standings contributions, and resubmits with corrected scores. Blocked if the next match in the bracket has already been completed. Admin only.
// @Tags         matches
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path string                    true "Tournament ID"
// @Param        matchId path string                    true "Match ID"
// @Param        body    body model.SubmitResultRequest true "Corrected set scores"
// @Success      200 {object} model.SuccessResponse{Data=object{result_detail=model.MatchDetail}, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Router       /tournaments/{id}/matches/{matchId}/result [put]
func (h *TournamentHandler) CorrectResult(c *gin.Context) {
	var req model.SubmitResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	detail, err := h.TournamentService.CorrectResult(c.Request.Context(), c.Param("id"), middleware.OrgID(c), c.Param("matchId"), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Result corrected successfully", Data: gin.H{"result_detail": detail}})
}

// ScheduleMatch godoc
// @Summary      Schedule match
// @Description  Sets the scheduled date and time for a match
// @Tags         matches
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path string                      true "Tournament ID"
// @Param        matchId path string                      true "Match ID"
// @Param        body    body model.ScheduleMatchRequest  true "Scheduled time (RFC3339)"
// @Success      200 {object} model.SuccessResponse{Data=object{scheduled_at=string}, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Router       /tournaments/{id}/matches/{matchId}/schedule [patch]
func (h *TournamentHandler) ScheduleMatch(c *gin.Context) {
	var req model.ScheduleMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	if err := h.TournamentService.ScheduleMatch(c.Param("matchId"), &req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Match scheduled successfully", Data: gin.H{"scheduled_at": req.ScheduledAt}})
}

// GetStageStandings godoc
// @Summary      Get stage standings
// @Description  Returns standings for all groups in a stage ordered by points. For group stages, use the group standings endpoint to get tiebreaker-resolved rankings per group.
// @Tags         standings
// @Produce      json
// @Security     BearerAuth
// @Param        id      path string true "Tournament ID"
// @Param        stageId path string true "Stage ID"
// @Success      200 {object} model.SuccessResponse{Data=object{Standings=[]model.Standings}, Message=string}
// @Failure      401 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /tournaments/{id}/stages/{stageId}/standings [get]
func (h *TournamentHandler) GetStageStandings(c *gin.Context) {
	standings, err := h.TournamentService.GetStageStandings(c.Param("stageId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Standings retrieved successfully", Data: gin.H{"Standings": standings}})
}

// GetGroupStandings godoc
// @Summary      Get group standings
// @Description  Returns standings for a single group with the full tiebreaker cascade applied (head-to-head → mutual set ratio → mutual game ratio → overall set ratio → overall game ratio → draw). Each entry includes a rank field.
// @Tags         standings
// @Produce      json
// @Security     BearerAuth
// @Param        id      path string true "Tournament ID"
// @Param        groupId path string true "Group ID"
// @Success      200 {object} model.SuccessResponse{Data=object{Standings=[]model.Standings}, Message=string}
// @Failure      401 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /tournaments/{id}/groups/{groupId}/standings [get]
func (h *TournamentHandler) GetGroupStandings(c *gin.Context) {
	standings, err := h.TournamentService.GetGroupStandings(c.Param("groupId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Group standings retrieved successfully", Data: gin.H{"Standings": standings}})
}
 