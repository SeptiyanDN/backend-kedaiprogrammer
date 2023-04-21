package categories

import (
	"time"
)

type Category struct {
	CategoryID   string    `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Tag          string    `json:"tag"`
	IsActive     bool      `json:"is_active"`
	ServiceID    string    `json:"service_id"`
	Created_at   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Updated_at   time.Time
}
