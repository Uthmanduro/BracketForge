package model

import (
	"encoding/json"
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
	ID              string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrganizationID  string          `gorm:"type:uuid;not null;index"                       json:"organization_id"`
	Name            string          `gorm:"not null"                                       json:"name"`
	Format          TournamentFormat `gorm:"not null"                                       json:"format"`
	Status          TournamentStatus `gorm:"not null;default:'draft'"                       json:"status"`
	BestOf          int             `gorm:"not null;default:3"                             json:"best_of"`
	ThirdPlaceMatch bool            `gorm:"not null;default:false"                         json:"third_place_match"`
	ScoringRules    json.RawMessage `gorm:"type:jsonb"                                     json:"scoring_rules,omitempty"`
	StartDate       *time.Time      `gorm:"default:null"                                   json:"start_date,omitempty"`
	EndDate         *time.Time      `gorm:"default:null"                                   json:"end_date,omitempty"`
	CreatedAt       time.Time       `gorm:"autoCreateTime"                                 json:"created_at"`
}

type TournamentRegistration struct {
	gorm.Model
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TournamentID string    `gorm:"type:uuid;not null;index"                       json:"tournament_id"`
	PlayerID     string    `gorm:"type:uuid;not null;index"                       json:"player_id"`
	Seed         *int      `gorm:"default:null"                                   json:"seed,omitempty"`
	Status       string    `gorm:"not null;default:'registered'"                  json:"status"`
	RegisteredAt time.Time `gorm:"autoCreateTime"                                 json:"registered_at"`
	// Joined
	PlayerName string `gorm:"-" json:"player_name,omitempty"`
}

func (TournamentRegistration) TableName() string { return "tournament_registrations" }

type CreateTournamentRequest struct {
	Name            string          `json:"name"              binding:"required"`
	Format          string          `json:"format"            binding:"required,oneof=single_elimination round_robin group_knockout"`
	BestOf          int             `json:"best_of"           binding:"required,oneof=3 5"`
	ThirdPlaceMatch bool            `json:"third_place_match"`
	ScoringRules    json.RawMessage `json:"scoring_rules"`
	StartDate       *time.Time      `json:"start_date"`
	EndDate         *time.Time      `json:"end_date"`
}

type RegisterPlayerRequest struct {
	PlayerID string `json:"player_id" binding:"required"`
	Seed     *int   `json:"seed"`
}

type UpdateRegistrationStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=registered checked_in withdrawn disqualified"`
}