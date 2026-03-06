package model

import "gorm.io/gorm"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleOrganizer  Role = "organizer"
	RoleGuest Role = "guest"
)

type User struct {
	gorm.Model
	OrganizationID uint `gorm:"not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Role Role `gorm:"not null;default:'user'" json:"role"`
}