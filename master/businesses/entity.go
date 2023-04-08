package businesses

import "time"

type Business struct {
	BusinessID          string    `json:"business_id"`
	BusinessName        string    `json:"business_name"`
	BusinessDescription string    `json:"business_description"`
	Domain              string    `json:"domain"`
	IsActive            bool      `json:"is_active,omitempty" gorm:"default:true"`
	Created_at          time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Updated_at          time.Time
}
