package businesses

import "time"

type Business struct {
	BusinessID    string    `json:"business_id"`
	Business_name string    `json:"business_name"`
	Domain        string    `json:"domain"`
	IsActive      bool      `json:"is_active"`
	Created_at    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Updated_at    time.Time
}
