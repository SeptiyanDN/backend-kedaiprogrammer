package users

import "time"

type User struct {
	ID         int       `json:"id"`
	Uuid       string    `json:"uuid"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Token      string    `json:"token"`
	Created_at time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	Updated_at time.Time
}
