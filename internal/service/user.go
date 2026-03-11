package service

import (
	"context"
	"fmt"

	"github.com/uthmanduro/BracketForge/internal/auth"
	"github.com/uthmanduro/BracketForge/internal/config"
	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/repository"
)

type AuthService struct {
	config *config.Config
	userRepo *repository.UserRepo
	// Add any dependencies needed for authentication, e.g. user repository, JWT secret, etc.
}

func NewAuthService(cfg *config.Config, userRepo *repository.UserRepo) *AuthService {
	return &AuthService{
		config: cfg,
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// Validate user credentials here, e.g. check username and password against the database
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("User not found, please create an account")
	}

	if !auth.VerifyPassword(password, user.Password) {
		return "", fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, string(user.Role), s.config.JWTSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Register(ctx context.Context, email, password string, organizationID uint) (*model.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("email already in use")
	}

	// Hash the password before storing it
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create new user in the database
	user := &model.User{
		Email: email,
		Password: hashedPassword,
		Role: "user", // Default role, can be changed as needed
		OrganizationID: organizationID,
	}

	s.userRepo.Create(ctx, user)
	return user, nil
}