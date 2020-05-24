package blog

import "go-blog/internal/db"

func List() ([]Blog, error) {
	db := db.GetConnection()

	results, err := db.Query("SELECT `id`, `title`, `content`, `description`, `slug`, `created_at` FROM `blogs`")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	blogs := make([]Blog, 0)

	for results.Next() {
		var blog Blog

		err = results.Scan(&blog.Id, &blog.Title, &blog.Content, &blog.Description, &blog.Slug, &blog.CreatedAt)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)

	}

	return blogs, nil
}
