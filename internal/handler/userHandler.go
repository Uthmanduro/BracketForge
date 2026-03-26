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

// Register godoc
// @Summary      Register user
// @Description  Creates a new user within the authenticated organisation. Admin only.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body model.RegisterUserRequest true "User details"
// @Success      201 {object} model.User
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Failure      403 {object} model.ErrorResponse
// @Router       /users/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req model.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.ErrorResponse{Error: err.Error()})
		return
	}
	email := req.Email
	password := req.Password
	role := req.Role
	

	// Basic validation
	if email == "" || password == "" {
		c.JSON(400, model.ErrorResponse{Error: "Email and password are required"})
		return
	}

	user, err := h.UserService.Register(middleware.OrgID(c), email, password, role)
	if err != nil {
		c.JSON(500, model.ErrorResponse{Error: "Failed to register user"})
		return
	}

	c.JSON(201, model.SuccessResponse{Message: "User registered successfully", Data: user})
}


// Login godoc
// @Summary      Login
// @Description  Authenticates with email and password. Returns a JWT token to use in the Authorization header.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body model.LoginRequest true "Credentials"
// @Success      200 {object} model.SuccessResponse{Data=map[string]string, Message=string}
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Router       /users/login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.ErrorResponse{Error: "Invalid request data"})
		return
	}
	email := req.Email
	password := req.Password

	if email == "" || password == "" {
		c.JSON(400, model.ErrorResponse{Error: "Email and password are required"})
		return
	}

	token, err := h.UserService.Login(email, password)
	if err != nil {
		c.JSON(401, model.ErrorResponse{Error: "Invalid email or password"})
		return
	}

	c.JSON(200, model.SuccessResponse{Message: "Login successful", Data: gin.H{"token": token}})
}

// Me godoc
// @Summary      Get current user
// @Description  Returns the profile of the currently authenticated user
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} model.SuccessResponse{Data=model.User, Message=string}
// @Failure      401 {object} model.ErrorResponse
// @Failure      404 {object} model.ErrorResponse
// @Router       /users/me [get]
func (h *UserHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.UserService.GetByID(userID.(string))
	if err != nil {
		c.JSON(404, model.ErrorResponse{Error: "User not found"})
		return
	}
	c.JSON(200, model.SuccessResponse{Message: "User retrieved successfully", Data: user})
}