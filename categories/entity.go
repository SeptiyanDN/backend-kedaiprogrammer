package categories

import (
	"time"
)

type Category struct {
	CategoryID   string    `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Slug         string    `json:"slug"`
	IsActive     bool      `json:"isActive"`
	BusinessID   string    `json:"business_id"`
	Created_at   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Updated_at   time.Time
}
