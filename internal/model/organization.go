package model

import "gorm.io/gorm"


type Organization struct {
	gorm.Model
	Name      string    `gorm:"not null" json:"name"`

	// Associations
	Users	 []User    `gorm:"foreignKey:OrganizationID" json:"users"`
	Tournaments []Tournament `gorm:"foreignKey:OrganizationID" json:"tournaments"`
}
