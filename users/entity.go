package users

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	Uuid       string         `json:"uuid"`
	Username   string         `json:"username"`
	Email      string         `json:"email"`
	Password   string         `json:"password"`
	Token      string         `json:"token"`
	UserInfo   postgres.Jsonb `gorm:"type:jsonb" json:"user_info"`
	BusinessID string         `json:"business_id"`
	Role       string         `json:"role,omitempty" gorm:"default:user"`
	IsActive   bool           `json:"is_active,omitempty" gorm:"default:true"`
	Created_at time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Updated_at time.Time
}
