package model

import "gorm.io/gorm"

type StageType string

const (
	GroupStage StageType = "group"
	KnockoutStage StageType = "knockout"
	RoundRobinStage StageType = "round_robin"
)

type Stage struct {
	gorm.Model
	TournamentID uint `gorm:"not null" json:"tournament_id"`
	Tournament   Tournament `gorm:"foreignKey:TournamentID" json:"tournament"`
	Name  string `gorm:"not null" json:"name"`
	Type  StageType `gorm:"not null" json:"type"`
	StageOrder int `gorm:"not null" json:"stage_order"`

	Matches []Match `gorm:"foreignKey:GroupID" json:"matches"`
	Groups  []Group `gorm:"foreignKey:StageID" json:"groups"`
}

