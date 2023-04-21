package articles

import "time"

type Article struct {
	ArticleID   string    `json:"article_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	Slug        string    `json:"slug"`
	CategoryID  string    `json:"category_id"`
	AuthorID    string    `json:"author_id"`
	MainImage   string    `json:"main_image"`
	PublisedAt  time.Time `json:"published_at,omitempty" gorm:"default:null;type:timestamp"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      int       `json:"status,omitempty" gorm:"default:0"`
}
