package handler

import (
	// "fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/service"
)

type UserHandler struct {
	AuthService *service.AuthService
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    OrgID     string `json:"organization_id" binding:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
    Email string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=8"`
}


func NewUserHandler(authService *service.AuthService) *UserHandler {
	return &UserHandler{
		AuthService: authService,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}
	email := req.Email
	password := req.Password
	orgId := req.OrgID

	// Basic validation
	if email == "" || password == "" || orgId == "" {
		c.JSON(400, gin.H{"error": "Email, password, and organization ID are required"})
		return
	}

	stringOrgId, err := strconv.Atoi(orgId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid organization ID"})
		return
	}

	user, err := h.AuthService.Register(c.Request.Context(), email, password, uint(stringOrgId))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(201, gin.H{"message": "User registered successfully", "user_id": user.ID})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req LoginUserRequest
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

	token, err := h.AuthService.Login(c.Request.Context(), email, password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(200, gin.H{"message": "Login successful", "token": token})
}