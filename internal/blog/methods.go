package blog

import (
	"database/sql"
	"go-blog/internal/db"
	"go-blog/pkg/utils"
)

func List(cursor int) ([]Blog, error) {
	db := db.GetConnection()

	var results *sql.Rows
	var err error

	queryString := "SELECT `id`, `title`, `content`, `description`, `slug`, `image`, `created_at` FROM `blogs` ORDER BY `id` DESC LIMIT 10"
	if cursor != 0 {
		queryString = "SELECT `id`, `title`, `content`, `description`, `slug`, `image`, `created_at` FROM `blogs` WHERE `id` < ? ORDER BY `id` DESC LIMIT 10"
		utils.Logg(queryString)
		results, err = db.Query(queryString, cursor)
	} else {
		results, err = db.Query(queryString)
	}

	if err != nil {
		return nil, err
	}
	defer results.Close()

	blogs := make([]Blog, 0)

	for results.Next() {
		var blog Blog

		err = results.Scan(&blog.Id, &blog.Title, &blog.Content, &blog.Description, &blog.Slug, &blog.Image, &blog.CreatedAt)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)

	}

	return blogs, nil
}
