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
	
	// idUint, err := strconv.ParseUint(id, 10, 64)
	// if err != nil {
	// 	c.JSON(400, gin.H{"error": "Invalid organization ID"})
	// 	return
	// }

	org, err := h.orgService.GetByID(id)
	if err != nil {
		c.JSON(404, model.ErrorResponse{Error: "Organization not found"})
		return
	}

	c.JSON(200, model.SuccessResponse{Message: "Organization retrieved successfully", Data: org})
}

// func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
// 	// Get organization ID from URL parameter and new name from form data
// 	id := c.Param("id")
// 	var req CreateUpdateOrganizationRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(400, gin.H{"error": "Invalid request data"})
// 		return
// 	}

// 	name := req.Name
// 	if name == "" {
// 		c.JSON(400, gin.H{"error": "Organization name is required"})
// 		return
// 	}

// 	// Convert ID to uint
// 	// idUint, err := strconv.ParseUint(id, 10, 64)
// 	// if err != nil {
// 	// 	c.JSON(400, gin.H{"error": "Invalid organization ID"})
// 	// 	return
// 	// }
	
// 	// Retrieve existing organization, update name, and save changes
// 	org, err := h.orgService.GetByID(id)
// 	if err != nil {
// 		c.JSON(404, gin.H{"error": "Organization not found"})
// 		return
// 	}
// 	org.Name = name

// 	// Save updated organization
// 	err = h.orgService.UpdateOrganization(c.Request.Context(), org)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": "Failed to update organization"})
// 		return
// 	}

// 	c.JSON(200, gin.H{"message": "Organization updated successfully", "organization": org})
// }

// func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
// 	id := c.Param("id")
// 	idUint, err := strconv.ParseUint(id, 10, 64)
// 	if err != nil {
// 		c.JSON(400, gin.H{"error": "Invalid organization ID"})
// 		return
// 	}

// 	err = h.orgService.DeleteOrganization(c.Request.Context(), uint(idUint))
// 	if err != nil {
// 		c.JSON(404, gin.H{"error": "Organization not found"})
// 		return
// 	}

// 	c.JSON(204, gin.H{"message": "Organization deleted successfully"})
// }