package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req CreateUpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}
	name := req.Name

	if name == "" {
		c.JSON(400, gin.H{"error": "Organization name is required"})
		return
	}
	org, err := h.orgService.CreateOrganization(c.Request.Context(), name)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create organization"})
		return
	}
	c.JSON(201, gin.H{"message": "Organization created successfully", "organization_id": org})
}

func (h *OrganizationHandler) GetOrganizationByID(c *gin.Context) {
	id := c.Param("id")
	
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid organization ID"})
		return
	}

	org, err := h.orgService.GetOrganizationByID(c.Request.Context(), uint(idUint))
	if err != nil {
		c.JSON(404, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(200, gin.H{"organization": org, "message": "Organization retrieved successfully"})
}

func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	// Get organization ID from URL parameter and new name from form data
	id := c.Param("id")
	var req CreateUpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	name := req.Name
	if name == "" {
		c.JSON(400, gin.H{"error": "Organization name is required"})
		return
	}

	// Convert ID to uint
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid organization ID"})
		return
	}
	
	// Retrieve existing organization, update name, and save changes
	org, err := h.orgService.GetOrganizationByID(c.Request.Context(), uint(idUint))
	if err != nil {
		c.JSON(404, gin.H{"error": "Organization not found"})
		return
	}
	org.Name = name

	// Save updated organization
	err = h.orgService.UpdateOrganization(c.Request.Context(), org)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update organization"})
		return
	}

	c.JSON(200, gin.H{"message": "Organization updated successfully", "organization": org})
}

func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid organization ID"})
		return
	}

	err = h.orgService.DeleteOrganization(c.Request.Context(), uint(idUint))
	if err != nil {
		c.JSON(404, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(204, gin.H{"message": "Organization deleted successfully"})
}