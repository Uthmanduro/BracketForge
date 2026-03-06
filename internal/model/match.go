package model

import "gorm.io/gorm"

type MatchStatus string

const (
	MatchStatusPending MatchStatus = "pending"
	MatchStatusCompleted MatchStatus = "completed"
)

type Match struct {
	gorm.Model
	StageID uint `gorm:"not null" json:"stage_id"`
	Stage   Stage `gorm:"foreignKey:StageID" json:"stage"`

	GroupID *uint `json:"group_id,omitempty"`
	Group   *Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`

	Player1ID uint `gorm:"not null" json:"player1_id"`
	Player1   Player `gorm:"foreignKey:Player1ID" json:"player1"`

	Player2ID uint `gorm:"not null" json:"player2_id"`
	Player2   Player `gorm:"foreignKey:Player2ID" json:"player2"`

	WinnerID *uint `json:"winner_id,omitempty"`
	Winner   *Player `gorm:"foreignKey:WinnerID" json:"winner,omitempty"`

	Status MatchStatus `gorm:"not null;default:'pending'" json:"status"`
	Round  int `gorm:"not null" json:"round"`

	Player1Score *int `json:"player1_score,omitempty"`
	Player2Score *int `json:"player2_score,omitempty"`
}