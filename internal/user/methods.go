package user

import (
	"database/sql"
	"go-blog/internal/db"
	"errors"
)

func List() ([]User, error) {
	db := db.GetConnection()

	results, err := db.Query("SELECT `id`, `name`, `username`, `email` FROM `users`")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	emps := make([]User, 0)

	for results.Next() {
		var emp User

		err = results.Scan(&emp.Id, &emp.Name, &emp.Username, &emp.Email)
		if err != nil {
			return nil, err
		}

		emps = append(emps, emp)

	}

	return emps, nil
}

func Read(id int64) (*User, error) {
	var emp User
	db := db.GetConnection()

	results := db.QueryRow("SELECT `id`, `name`, `username`, `email`, `role` FROM `users` WHERE `id` = ? ", id)
	err := results.Scan(&emp.Id, &emp.Name, &emp.Username, &emp.Email)
	if err == sql.ErrNoRows {
		return nil, errors.New("Invalid id.")
	} else if err != nil {
		return nil, err
	}

	return &emp, nil
}

func Create(user User) (int64, error) {
	db := db.GetConnection()

	results, err := db.Exec("INSERT INTO `users` (`name`, `username`, `email`, `role`, `password`) VALUES (?,?,?,?,?)", user.Name, user.Username, user.Email, user.Role, user.Password)
	if err != nil {
		return 0, err
	}

	lid, err := results.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lid, nil
}

func Update(user User) error {
	db := db.GetConnection()

	_, err := db.Exec("UPDATE `users` SET `name` = ?, `username` = ?, `email` = ?, `role` = ?, `password` = ? WHERE `id` = ?", user.Name, user.Username, user.Email, user.Role, user.Password, user.Id)
	return err
}

func Delete(id int64) error {
	db := db.GetConnection()

	_, err := db.Exec("DELETE FROM `users` WHERE `id` = ?", id)
	return err
}