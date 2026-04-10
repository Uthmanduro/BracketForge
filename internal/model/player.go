package model

import (
	"encoding/json"
	"time"

)


type Player struct {
	ID             string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrganizationID string          `gorm:"type:uuid;not null;index"                       json:"organization_id"`
	Name           string          `gorm:"not null"                                       json:"name"`
	Email          *string         `gorm:"default:null"                                   json:"email,omitempty"`
	Metadata       json.RawMessage `gorm:"type:jsonb"                                     json:"metadata,omitempty"`
	CreatedAt      time.Time       `gorm:"autoCreateTime"                                 json:"created_at"`
}
 
type CreatePlayerRequest struct {
	Name     string          `json:"name"     binding:"required"`
	Email    *string         `json:"email"`
	Metadata json.RawMessage `json:"metadata"`
}