package db

import "database/sql"

type User struct {
	ID       int
	Username string
	Password string
	Hash     string
}

func (u User) Label() string {
	return u.Username
}

func LoadUsers(db *sql.DB) ([]User, error) {

	rows, err := db.Query(
		"SELECT id, username, password, hash FROM users",
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []User

	for rows.Next() {
		var u User

		err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Hash)
		if err != nil {
			return nil, err
		}

		list = append(list, u)
	}

	return list, nil
}

func CreateUser(db *sql.DB, u User) error {

	_, err := db.Exec(
		"INSERT INTO users(username, password, hash) VALUES (?, ?, ?)",
		u.Username,
		u.Password,
		u.Hash,
	)

	return err
}

func UpdateUser(db *sql.DB, u User) error {

	_, err := db.Exec(
		`UPDATE users 
		 SET username = ?, password = ?, hash = ?
		 WHERE id = ?`,
		u.Username,
		u.Password,
		u.Hash,
		u.ID,
	)

	return err
}

func DeleteUser(db *sql.DB, id int) error {

	_, err := db.Exec(
		`DELETE FROM users WHERE id = ?`,
		id,
	)

	return err
}
