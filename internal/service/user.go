package service

import (
	"errors"
	"fmt"

	"github.com/uthmanduro/BracketForge/internal/auth"
	"github.com/uthmanduro/BracketForge/internal/model"
	"github.com/uthmanduro/BracketForge/internal/repository"
)


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
 
func (s *UserService) GetByEmail(email string) (*model.User, error) {
	return s.repo.GetByEmail(email)
}