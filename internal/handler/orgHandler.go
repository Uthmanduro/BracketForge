package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/service"
)

type OrganizationHandler struct {
	orgService *service.OrganizationService
}



type CreateUpdateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}

func NewOrganizationHandler(orgService *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
	}
}

// Create godoc
// @Summary      Create organisation
// @Description  Creates a new organisation. This is the onboarding entry point — no auth required.
// @Tags         organisations
// @Accept       json
// @Produce      json
// @Param        body body CreateUpdateOrganizationRequest true "Organisation name"
// @Success      201 {object} model.SuccessResponse{Data=model.Organization, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /organizations [post]
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req CreateUpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.ErrorResponse{Error: err.Error()})
		return
	}
	name := req.Name

	if name == "" {
		c.JSON(400, model.ErrorResponse{Error: "Organization name is required"})
		return
	}

	checkName, err := h.orgService.GetByName(name)
	if err == nil && checkName != nil {
		c.JSON(400, model.ErrorResponse{Error: "Organization name already exists"})
		return
	}

	org, err := h.orgService.Create(name)
	if err != nil {
		c.JSON(500, model.ErrorResponse{Error: "Failed to create organization"})
		return
	}
	c.JSON(201, model.SuccessResponse{Message: "Organization created successfully", Data: org})
}

// GetByID godoc
// @Summary      Get organisation
// @Description  Returns the organisation for the authenticated user's org ID
// @Tags         organisations
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Organisation ID"
// @Success      200 {object} model.SuccessResponse{Data=model.Organization, Message=string}
// @Failure      401 {object} model.ErrorResponse
// @Failure      404 {object} model.ErrorResponse
// @Router       /organizations/{id} [get]
func (h *OrganizationHandler) GetOrganizationByID(c *gin.Context) {
	id := c.Param("id")

	org, err := h.orgService.GetByID(id)
	if err != nil {
		c.JSON(404, model.ErrorResponse{Error: "Organization not found"})
		return
	}

	c.JSON(200, model.SuccessResponse{Message: "Organization retrieved successfully", Data: org})
}

// GetAll godoc
// @Summary      List organisations
// @Description  Returns a list of all organisations. This is primarily for testing and debugging purposes.
// @Tags         organisations
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} model.SuccessResponse{Data=[]model.Organization, Message=string}
// @Failure      401 {object} model.ErrorResponse
// @Failure      500 {object} model.ErrorResponse
// @Router       /organizations [get]
func (h *OrganizationHandler) GetAllOrganization(c *gin.Context) {
	orgs, err := h.orgService.GetAll()
	if err != nil {
		c.JSON(500, model.ErrorResponse{Error: "Failed to retrieve organizations"})
		return
	}

	c.JSON(200, model.SuccessResponse{Message: "Organizations retrieved successfully", Data: orgs})
}