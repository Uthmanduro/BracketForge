package model

import (
	"time"

	"gorm.io/gorm"
)

type MatchStatus string

const (
	MatchStatusPending MatchStatus = "pending"
	MatchStatusCompleted MatchStatus = "completed"
)

type Match struct {
	gorm.Model
	ID                   string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StageID              string     `gorm:"type:uuid;not null;index"                       json:"stage_id"`
	GroupID              *string    `gorm:"type:uuid;index;default:null"                   json:"group_id,omitempty"`
	NextMatchID          *string    `gorm:"type:uuid;default:null"                         json:"next_match_id,omitempty"`
	NextMatchSlot        *int       `gorm:"default:null"                                   json:"next_match_slot,omitempty"`
	LoserNextMatchID     *string    `gorm:"type:uuid;default:null"                         json:"loser_next_match_id,omitempty"`
	LoserNextMatchSlot   *int       `gorm:"default:null"                                   json:"loser_next_match_slot,omitempty"`
	Round                *int       `gorm:"default:null"                                   json:"round,omitempty"`
	MatchPosition        *int       `gorm:"default:null"                                   json:"match_position,omitempty"`
	BestOf               int        `gorm:"not null;default:3"                             json:"best_of"`
	IsBye                bool       `gorm:"not null;default:false"                         json:"is_bye"`
	WinnerRegistrationID *string    `gorm:"type:uuid;default:null"                         json:"winner_registration_id,omitempty"`
	Status               MatchStatus `gorm:"not null;default:'pending'"                     json:"status"`
	ScheduledAt          *time.Time `gorm:"default:null"                                   json:"scheduled_at,omitempty"`
	CompletedAt          *time.Time `gorm:"default:null"                                   json:"completed_at,omitempty"`
	CreatedAt            time.Time  `gorm:"autoCreateTime"                                 json:"created_at"`
}

type MatchParticipant struct {
	gorm.Model
	ID             string  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	MatchID        string  `gorm:"type:uuid;not null;index"                       json:"match_id"`
	RegistrationID string  `gorm:"type:uuid;not null;index"                       json:"registration_id"`
	Slot           int     `gorm:"not null"                                       json:"slot"`
	Result         *string `gorm:"default:null"                                   json:"result,omitempty"`
	// Joined
	PlayerName string `gorm:"-" json:"player_name,omitempty"`
}

func (MatchParticipant) TableName() string { return "match_participants" }

type SetScore struct {
	gorm.Model
	ID               string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	MatchID          string `gorm:"type:uuid;not null;index"                       json:"match_id"`
	SetNumber        int    `gorm:"not null"                                       json:"set_number"`
	P1Games          int    `gorm:"not null"                                       json:"p1_games"`
	P2Games          int    `gorm:"not null"                                       json:"p2_games"`
	IsTiebreak       bool   `gorm:"not null;default:false"                         json:"is_tiebreak"`
	P1TiebreakPoints *int   `gorm:"default:null"                                   json:"p1_tiebreak_points,omitempty"`
	P2TiebreakPoints *int   `gorm:"default:null"                                   json:"p2_tiebreak_points,omitempty"`
}

func (SetScore) TableName() string { return "set_scores" }

type Standings struct {
	gorm.Model
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StageID        string    `gorm:"type:uuid;not null;index"                       json:"stage_id"`
	GroupID        *string   `gorm:"type:uuid;index;default:null"                   json:"group_id,omitempty"`
	RegistrationID string    `gorm:"type:uuid;not null;index"                       json:"registration_id"`
	Played         int       `gorm:"not null;default:0"                             json:"played"`
	Wins           int       `gorm:"not null;default:0"                             json:"wins"`
	Losses         int       `gorm:"not null;default:0"                             json:"losses"`
	Points         int       `gorm:"not null;default:0"                             json:"points"`
	SetsWon        int       `gorm:"not null;default:0"                             json:"sets_won"`
	SetsLost       int       `gorm:"not null;default:0"                             json:"sets_lost"`
	GamesWon       int       `gorm:"not null;default:0"                             json:"games_won"`
	GamesLost      int       `gorm:"not null;default:0"                             json:"games_lost"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"                                 json:"updated_at"`
	// Joined
	PlayerName string `gorm:"-" json:"player_name,omitempty"`
	Rank       int    `gorm:"-" json:"rank"`
}

type SetScoreInput struct {
	SetNumber        int  `json:"set_number"         binding:"required"`
	P1Games          int  `json:"p1_games"           binding:"min=0"`
	P2Games          int  `json:"p2_games"           binding:"min=0"`
	IsTiebreak       bool `json:"is_tiebreak"`
	P1TiebreakPoints *int `json:"p1_tiebreak_points"`
	P2TiebreakPoints *int `json:"p2_tiebreak_points"`
}

type SubmitResultRequest struct {
	Sets []SetScoreInput `json:"sets" binding:"required,min=1"`
}

type MatchDetail struct {
	Match        *Match              `json:"match"`
	Participants []*MatchParticipant `json:"participants"`
	Sets         []*SetScore         `json:"sets"`
}

type ScheduleMatchRequest struct {
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
}

// MutualSetStats holds per-registration mutual-match tallies for tiebreaking.
type MutualSetStats struct {
	SetsWon   int
	SetsLost  int
	GamesWon  int
	GamesLost int
}