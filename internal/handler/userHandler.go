package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/middleware"
	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/service"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req model.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	email := req.Email
	password := req.Password
	role := req.Role
	

	// Basic validation
	if email == "" || password == "" {
		c.JSON(400, gin.H{"error": "Email and password are required"})
		return
	}

	user, err := h.UserService.Register(middleware.OrgID(c), email, password, role)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(201, gin.H{"message": "User registered successfully", "user_id": user.ID})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}
	email := req.Email
	password := req.Password

	if email == "" || password == "" {
		c.JSON(400, gin.H{"error": "Email and password are required"})
		return
	}

	token, err := h.UserService.Login(email, password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(200, gin.H{"message": "Login successful", "token": token})
}

func (h *UserHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.UserService.GetByID(userID.(string))
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, user)
}