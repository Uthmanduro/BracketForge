package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	TournamentID uint `gorm:"not null" json:"tournament_id"`
	Tournament   Tournament `gorm:"foreignKey:TournamentID" json:"tournament"`
	StageID uint `gorm:"not null" json:"stage_id"`
	Stage   Stage `gorm:"foreignKey:StageID" json:"stage"`
	Name  string `gorm:"not null" json:"name"`


	Matches []Match `gorm:"foreignKey:GroupID" json:"matches"`
}