package db

import "database/sql"

type Host struct {
	ID     int
	Name   string
	IP     string
	Domain string
	Role   string
}

func LoadHosts(db *sql.DB) ([]Host, error) {

	rows, err := db.Query(
		"SELECT id, name, ip, domain, role FROM machines",
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Host

	for rows.Next() {
		var h Host

		err := rows.Scan(&h.ID, &h.Name, &h.IP, &h.Domain, &h.Role)
		if err != nil {
			return nil, err
		}

		list = append(list, h)
	}

	return list, nil
}

func InsertHost(db *sql.DB, h Host) error {

	_, err := db.Exec(
		"INSERT INTO machines(name, ip, domain, role) VALUES (?, ?, ?, ?)",
		h.Name,
		h.IP,
		h.Domain,
		h.Role,
	)

	return err
}

func UpdateHost(db *sql.DB, h Host) error {

	_, err := db.Exec(
		`UPDATE machines 
		 SET name = ?, ip = ?, domain = ?, role = ?
		 WHERE id = ?`,
		h.Name,
		h.IP,
		h.Domain,
		h.Role,
		h.ID,
	)

	return err
}

func DeleteHost(db *sql.DB, id int) error {

	_, err := db.Exec(
		"DELETE FROM machines WHERE id = ?",
		id,
	)

	return err
}
