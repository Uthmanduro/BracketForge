package model

import (
	"time"

)

type StageType string

const (
	GroupStage StageType = "group"
	KnockoutStage StageType = "knockout"
	RoundRobinStage StageType = "round_robin"
)

type Stage struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TournamentID string    `gorm:"type:uuid;not null;index"                       json:"tournament_id"`
	Name         string    `gorm:"not null"                                       json:"name"`
	Type         StageType `gorm:"not null"                                       json:"type"`
	StageOrder   int       `gorm:"not null"                                       json:"stage_order"`
	AdvanceCount *int      `gorm:"default:null"                                   json:"advance_count,omitempty"`
	BestOf       *int      `gorm:"default:null"                                   json:"best_of,omitempty"`
	CreatedAt    time.Time `gorm:"autoCreateTime"                                 json:"created_at"`
}

type Group struct {
	ID      string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StageID string `gorm:"type:uuid;not null;index"                       json:"stage_id"`
	Name    string `gorm:"not null"                                       json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type GroupRegistration struct {
	ID             string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	GroupID        string `gorm:"type:uuid;not null;index"                       json:"group_id"`
	RegistrationID string `gorm:"type:uuid;not null;index"                       json:"registration_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (GroupRegistration) TableName() string { return "group_registrations" }

type CreateStageRequest struct {
	Name         string `json:"name"          binding:"required"`
	Type         string `json:"type"          binding:"required,oneof=group knockout round_robin"`
	StageOrder   int    `json:"stage_order"   binding:"required"`
	AdvanceCount *int   `json:"advance_count"`
	BestOf       *int   `json:"best_of"`

}

type DrawRequest struct {
	NumberOfGroups int `json:"number_of_groups" binding:"required,min=2"`
}