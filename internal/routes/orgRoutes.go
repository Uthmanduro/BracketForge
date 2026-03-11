package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/handler"
)

func RegisterOrgRoutes(r *gin.RouterGroup, handler *handler.OrganizationHandler) {
	orgGroup := r.Group("/orgs")
	{
		orgGroup.POST("/", handler.CreateOrganization)
		orgGroup.GET("/:id", handler.GetOrganizationByID)
		orgGroup.PUT("/:id", handler.UpdateOrganization)
		orgGroup.DELETE("/:id", handler.DeleteOrganization)
	}
}