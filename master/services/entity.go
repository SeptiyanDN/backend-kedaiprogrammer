package services

import (
	"time"
)

type Service struct {
	ServiceID   string    `json:"service_id"`
	ServiceName string    `json:"service_name"`
	IsActive    bool      `json:"is_active"`
	BusinessID  string    `json:"business_id"`
	Created_at  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Updated_at  time.Time
}
