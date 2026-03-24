package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/config"
	"github.com/uthmanduro/BracketForge/internal/handler"
	"github.com/uthmanduro/BracketForge/internal/middleware"
)

func RegisterOrgRoutes(r *gin.RouterGroup, cfg *config.Config, handler *handler.OrganizationHandler) {
	orgGroup := r.Group("/orgs")
	{
		orgGroup.POST("/", handler.CreateOrganization)
		orgGroup.Use(middleware.AuthMiddleware(cfg))
		orgGroup.GET("/:id", handler.GetOrganizationByID)
	}
}