package service

import (
	"errors"
	"fmt"

	"github.com/uthmanduro/BracketForge/internal/auth"
	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/repository"
)

// type AuthService struct {
// 	config *config.Config
// 	userRepo *repository.UserRepo
// 	// Add any dependencies needed for authentication, e.g. user repository, JWT secret, etc.
// }

// func NewAuthService(cfg *config.Config, userRepo *repository.UserRepo) *AuthService {
// 	return &AuthService{
// 		config: cfg,
// 		userRepo: userRepo,
// 	}
// }

// func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
// 	// Validate user credentials here, e.g. check username and password against the database
// 	user, err := s.userRepo.FindByEmail(ctx, email)
// 	if err != nil {
// 		return "", err
// 	}
// 	if user == nil {
// 		return "", fmt.Errorf("User not found, please create an account")
// 	}

// 	if !auth.VerifyPassword(password, user.Password) {
// 		return "", fmt.Errorf("invalid email or password")
// 	}

// 	// Generate JWT token
// 	token, err := auth.GenerateJWT(user.ID, string(user.Role), s.config.JWTSecret)
// 	if err != nil {
// 		return "", err
// 	}

// 	return token, nil
// }

// func (s *AuthService) Register(ctx context.Context, email, password string, organizationID uint) (*model.User, error) {
// 	// Check if user already exists
// 	existingUser, err := s.userRepo.FindByEmail(ctx, email)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if existingUser != nil {
// 		return nil, fmt.Errorf("email already in use")
// 	}

// 	// Hash the password before storing it
// 	hashedPassword, err := auth.HashPassword(password)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create new user in the database
// 	user := &model.User{
// 		Email: email,
// 		Password: hashedPassword,
// 		Role: "user", // Default role, can be changed as needed
// 		OrganizationID: organizationID,
// 	}

// 	s.userRepo.Create(ctx, user)
// 	return user, nil
// }


type UserService struct {
	repo      *repository.UserRepo
	jwtSecret string
}
 
func NewUserService(r *repository.UserRepo, secret string) *UserService {
	return &UserService{repo: r, jwtSecret: secret}
}
 
func (s *UserService) Register(orgID, email, password, role string) (*model.User, error) {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	u := &model.User{OrganizationID: orgID, Email: email, PasswordHash: hash, Role: model.Role(role)}
	return u, s.repo.Create(u)
}
 
func (s *UserService) Login(email, password string) (*model.LoginResponse, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if !auth.VerifyPassword(password, u.PasswordHash) {
		return nil, errors.New("password not match")
	}
	token, err := auth.GenerateJWT(u.ID,  string(u.Role), u.OrganizationID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	return &model.LoginResponse{Token: token, User: u}, nil
}
 
func (s *UserService) GetByID(id string) (*model.User, error) {
	return s.repo.GetByID(id)
}
 