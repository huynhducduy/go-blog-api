package blog

type Blog struct {
	Id          *int64 `json:"id,omitempty"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Slug        *string `json:"slug"`
	Image		*string `json:"image"`
	Content     *string `json:"content"`
	CreatedAt   *int64  `json:"created_at"`
}