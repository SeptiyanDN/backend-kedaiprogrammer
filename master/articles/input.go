package articles

type CreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
	CategoryID  string `json:"category_id"`
	AuthorID    string `json:"author_id"`
	MainImage   string `json:"main_image"`
	PublisedAt  string `json:"published_at"`
}
