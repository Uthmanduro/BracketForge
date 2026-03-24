package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/config"
	"github.com/uthmanduro/BracketForge/internal/handler"
	"github.com/uthmanduro/BracketForge/internal/middleware"
	// "github.com/uthmanduro/BracketForge/internal/middleware"clear
)

func RegisterUserRoutes(r *gin.RouterGroup, handler *handler.UserHandler, config *config.Config) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/login", handler.LoginUser)
		authGroup := userGroup.Group("/")
		authGroup.Use(middleware.AuthMiddleware(config))
		authGroup.POST("/register", middleware.RequireRole("admin"), handler.RegisterUser)
		authGroup.GET("/me", handler.Me)
	}
}