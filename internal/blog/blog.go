package blog

import "go-blog/internal/user"

type Blog struct {
	Id          *int64 `json:"id,omitempty"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Slug        *string `json:"slug"`
	Image		*string `json:"image"`
	Content     *string `json:"content"`
	CreatedAt   *int64  `json:"created_at"`
	UserId      *int64 `json:"user_id"`
	Tags 		*string `json:"tags"`
}

type BlogExtra struct {
	Blog
	User user.User `json:"user"`
}

type BlogFilter struct {
	Title *string
	UserId int64
	Tags *string
}

type BlogSortMethod string

const (
	SortByCreatedAtAsc = "created_at_asc"
	SortByCreatedAtDesc = "created_at_desc"
)