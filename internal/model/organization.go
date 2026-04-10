package model

import (
	"time"
)


type Organization struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"not null"                                        json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime"                                  json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Data	any    `json:"data,omitempty"`
}