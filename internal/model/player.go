package model

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	TournamentID uint `gorm:"not null" json:"tournament_id"`
	Tournament   Tournament `gorm:"foreignKey:TournamentID" json:"tournament"`
	Name  string `gorm:"not null" json:"name"`
	Seed  int    `json:"seed"`
	Ranking *int    `json:"ranking"`

	Matches []Match `gorm:"many2many:match_players;" json:"matches"`
}