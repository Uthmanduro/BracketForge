package model

import (
	"time"

	"gorm.io/gorm"
)

type TournamentFormat string
type TournamentStatus string

const (
	SingleElimination TournamentFormat = "single_elimination"
	RoundRobin TournamentFormat = "round_robin"
	GroupKnockout TournamentFormat = "group_knockout"
)

const (
	Draft TournamentStatus = "draft"
	Active TournamentStatus = "active"
	Completed TournamentStatus = "completed"
)



type Tournament struct {
	gorm.Model
	OrganizationID uint `gorm:"not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	Name        string           `gorm:"not null" json:"name"`
	Description string           `json:"description"`
	Format      TournamentFormat `gorm:"not null" json:"format"`
	Status      TournamentStatus `gorm:"not null;default:'draft'" json:"status"`
	StartDate   *time.Time       `json:"start_date"`
	EndDate     *time.Time       `json:"end_date"`

	Players []Player `gorm:"foreignKey:TournamentID" json:"players"`
	Stages  []Stage  `gorm:"foreignKey:TournamentID" json:"stages"`
}