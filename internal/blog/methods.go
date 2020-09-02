package blog

import (
	"database/sql"
	"go-blog/internal/config"
	"go-blog/internal/db"
	"go-blog/internal/tag"
	"go-blog/pkg/utils"
	"strings"
)


// Return extra
func List(cursor int, filter BlogFilter, sortMethod BlogSortMethod) ([]BlogExtra, error) {
	db := db.GetConnection()

	var results *sql.Rows
	var err error

	queryString := "SELECT `b`.`id`, `b`.`title`, `b`.`content`, `b`.`description`, `b`.`slug`, `b`.`image`, `b`.`created_at`, `b`.`user_id`, `b`.`tags`, `u`.`id`, `u`.`username`, `u`.`email`, `u`.`role`, `u`.`name`  FROM `blogs` as b, `users` as u WHERE `b`.`user_id` = `u`.`id`";
	variables := make([]interface{}, 0)

	if filter.Title != nil {
		queryString += " AND `b`.`title` LIKE ?"
		queryVar := "%" + *filter.Title + "%"
		variables = append(variables, queryVar)
	}

	if filter.UserId != 0 {
		queryString += " AND `b`.`user_id` = ?"
		variables = append(variables, filter.UserId)
	}

	if filter.Tags != nil {
		for _, v := range strings.Split(*filter.Tags, ",") {
			if v != "" {
				queryString += " AND `b`.`tags` LIKE ?"
				variables = append(variables, "%," + v + ",%")
			}
		}
	}

	if cursor != 0 {
		queryString += " AND `b`.`id` < ?"
		variables = append(variables, cursor)
	}

	switch sortMethod {
	case SortByCreatedAtAsc:
		queryString += " ORDER BY `b`.`created_at` ASC"
	case SortByCreatedAtDesc:
		queryString += " ORDER BY `b`.`created_at` DESC"
	default:
		queryString += " ORDER BY `b`.`id` DESC"
	}

	queryString += " LIMIT ?"
	variables = append(variables,config.GetConfig().ITEMS_PER_PAGE)

	// DONE evaluating query string

	utils.Logg(queryString)

	results, err = db.Query(queryString, variables...)

	if err != nil {
		return nil, err
	}
	defer results.Close()

	blogs := make([]BlogExtra, 0)

	for results.Next() {
		var blog BlogExtra

		err = results.Scan(&blog.Id, &blog.Title, &blog.Content, &blog.Description, &blog.Slug, &blog.Image, &blog.CreatedAt, &blog.UserId, &blog.Tags, &blog.User.Id, &blog.User.Username, &blog.User.Email, &blog.User.Role, &blog.User.Name)
		if err != nil {
			return nil, err
		}

		if blog.Tags != nil {
			var newTags = strings.Trim(*blog.Tags, ",")
			blog.Tags = &newTags
		}

		blogs = append(blogs, blog)

	}

	return blogs, nil
}

// Return extra
func Read(id int) (*BlogExtra, error) {
	db := db.GetConnection()

	var blog BlogExtra

	result := db.QueryRow("SELECT `b`.`id`, `b`.`title`, `b`.`content`, `b`.`description`, `b`.`slug`, `b`.`image`, `b`.`created_at`, `b`.`user_id`, `b`.`tags`, `u`.`id`, `u`.`username`, `u`.`email`, `u`.`role`, `u`.`name`  FROM `blogs` as b, `users` as u WHERE `b`.`user_id` = `u`.`id` AND `b`.`id` = ?", id)
	err := result.Scan(&blog.Id, &blog.Title, &blog.Content, &blog.Description, &blog.Slug, &blog.Image, &blog.CreatedAt, &blog.UserId, &blog.Tags, &blog.User.Id, &blog.User.Username, &blog.User.Email, &blog.User.Role, &blog.User.Name)
	if err != nil {
		return nil, err
	}

	if blog.Tags != nil {
		var newTags = strings.Trim(*blog.Tags, ",")
		blog.Tags = &newTags
	}

	return &blog, err
}

func Create(blog Blog) (int64, error) {
	db := db.GetConnection()

	if blog.Tags != nil {
		var newTags = "," + strings.Trim(*blog.Tags,",") + ","
		blog.Tags = &newTags
	}

	results, err := db.Exec("INSERT INTO `blogs` (`title`, `content`, `description`, `slug`, `image`, `created_at`, `user_id`, `tags`) VALUES (?,?,?,?,?,?,?,?)", blog.Title, blog.Content, blog.Description, blog.Slug, blog.Image, blog.CreatedAt, blog.UserId, blog.Tags)
	if err != nil {
		return 0, err
	}

	if blog.Tags != nil {
		for _, v := range strings.Split(*blog.Tags, ",") {
			if v != "" {
				var newTag tag.Tag
				newTag.Tag = &v
				tag.Create(newTag)
			}
		}
	}

	lid, err := results.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lid, nil
}

func Update(blog Blog) error {
	db := db.GetConnection()

	var newTags = "," + strings.Trim(*blog.Tags,",") + ","
	blog.Tags = &newTags

	_, err := db.Exec("UPDATE `blogs` SET `title` = ?, `description` = ?, `slug` = ?, `image` = ?, `content` = ?, `tags` = ? WHERE `id` = ?", blog.Title, blog.Description, blog.Slug, blog.Image, blog.Content, blog.Tags, blog.Id)
	return err
}

func Delete(id int64) error {
	db := db.GetConnection()

	_, err := db.Exec("DELETE FROM `blogs` WHERE `id` = ?", id)
	return err
}
