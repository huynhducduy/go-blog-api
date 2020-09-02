package tag

import (
	"database/sql"
	"go-blog/internal/db"
)

func List() ([]Tag, error) {
	db := db.GetConnection()

	results, err := db.Query("SELECT `tag`, `title`, `description` FROM `tags`")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	emps := make([]Tag, 0)

	for results.Next() {
		var emp Tag

		err = results.Scan(&emp.Tag, &emp.Title, &emp.Description)
		if err != nil {
			return nil, err
		}

		emps = append(emps, emp)

	}

	return emps, nil
}

func Read(tag string) (*Tag, error) {
	var emp Tag
	db := db.GetConnection()

	results := db.QueryRow("SELECT `tag`, `title`, `description` FROM `tags` WHERE `tag` = ?", tag)
	err := results.Scan(&emp.Tag, &emp.Title, &emp.Description)
	if err == sql.ErrNoRows {
		emp.Tag = &tag
		return &emp, nil
	} else if err != nil {
		return nil, err
	}

	return &emp, nil
}

// Create or not -> using when create blog
func Create(tag Tag) (error) {
	db := db.GetConnection()

	_, err := db.Exec("INSERT INTO `tags` (`tag`, `title`, `description`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE tag = ?", tag.Tag, tag.Title, tag.Description, tag.Tag)
	if err != nil {
		return err
	}

	return nil
}

// Update or create
func Update(tag Tag) error {
	db := db.GetConnection()

	_, err := db.Exec("INSERT INTO `tags` (`tag`, `title`, `description`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE `title` = ?, `description` = ?", tag.Tag, tag.Title, tag.Description, tag.Title, tag.Description)
	return err
}

func Delete(tag string) error {
	db := db.GetConnection()

	_, err := db.Exec("DELETE FROM `tags` WHERE `tag` = ?", tag)
	return err
}