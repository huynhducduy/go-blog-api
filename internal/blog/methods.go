package blog

import (
	"database/sql"
	"errors"
	"go-blog/internal/db"
	"go-blog/pkg/utils"
)

func List(cursor int) ([]Blog, error) {
	db := db.GetConnection()

	var results *sql.Rows
	var err error

	queryString := "SELECT `id`, `title`, `content`, `description`, `slug`, `image`, `created_at` FROM `blogs` ORDER BY `id` DESC LIMIT 5"
	if cursor != 0 {
		queryString = "SELECT `id`, `title`, `content`, `description`, `slug`, `image`, `created_at` FROM `blogs` WHERE `id` < ? ORDER BY `id` DESC LIMIT 5"
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

func Read(id int) (*Blog, error) {
	db := db.GetConnection()

	var blog Blog

	result := db.QueryRow("SELECT `id`, `title`, `content`, `description`, `slug`, `image`, `created_at` FROM `blogs` WHERE `id` = ?", id)
	err := result.Scan(&blog.Id, &blog.Title, &blog.Content, &blog.Description, &blog.Slug, &blog.Image, &blog.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("Invalid id.")
	} else if err != nil {
		return nil, err
	}

	return &blog, err
}

func Create(blog Blog) (int64, error) {
	db := db.GetConnection()

	results, err := db.Exec("INSERT INTO `blogs` (`title`, `content`, `description`, `slug`, `image`, `created_at`) VALUES (?,?,?,?,?,?)", blog.Title, blog.Content, blog.Description, blog.Slug, blog.Image, blog.CreatedAt)
	if err != nil {
		return 0, err
	}

	lid, err := results.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lid, nil
}

func Update(blog Blog) error {
	db := db.GetConnection()

	_, err := db.Exec("UPDATE `blogs` SET `title` = ?, `description` = ?, `slug` = ?, `image` = ?, `content` = ?, `created_at` = ?", blog.Title, blog.Description, blog.Slug, blog.Image, blog.Content, blog.CreatedAt)
	return err
}

func Delete(id int64) error {
	db := db.GetConnection()

	_, err := db.Exec("DELETE FROM `blogs` WHERE `ID` = ?", id)
	return err
}
