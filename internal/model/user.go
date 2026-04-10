package model

import (
	"time"

)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleOrganizer  Role = "organizer"
	RoleGuest Role = "guest"
)

type User struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrganizationID string    `gorm:"type:uuid;not null;index"                       json:"organization_id"`
	Email          string    `gorm:"uniqueIndex;not null"                           json:"email"`
	PasswordHash   string    `gorm:"not null"                                       json:"-"`
	Role Role `gorm:"not null;default:'organizer'" json:"role"`
	CreatedAt      time.Time `gorm:"autoCreateTime"                                 json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type RegisterUserRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role"     binding:"required,oneof=admin organizer viewer"`
	OrgID	string `json:"org_id"   binding:"required,uuid"`
}