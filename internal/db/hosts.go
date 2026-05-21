package db

import "database/sql"

type Subdomain struct {
	ID        int
	HostID    int
	Subdomain string
}

type Host struct {
	ID         int
	Name       string
	IP         string
	Domain     string
	Role       string
	Subdomains []Subdomain
}

func (h Host) Label() string {
	return h.Name
}

func LoadHosts(db *sql.DB) ([]Host, error) {

	rows, err := db.Query(
		"SELECT id, name, ip, domain, role FROM hosts",
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

	// load subdomains
	for i := range list {

		subRows, err := db.Query(`
			SELECT id, subdomain 
			FROM subdomains
			WHERE host_id = ?
		`, list[i].ID)

		if err != nil {
			return nil, err
		}

		for subRows.Next() {
			var s Subdomain
			s.HostID = list[i].ID

			if err := subRows.Scan(&s.ID, &s.Subdomain); err != nil {
				subRows.Close()
				return nil, err
			}

			list[i].Subdomains = append(list[i].Subdomains, s)
		}

		subRows.Close()
	}

	return list, nil
}

func CreateHost(db *sql.DB, h Host) error {

	_, err := db.Exec(
		"INSERT INTO hosts(name, ip, domain, role) VALUES (?, ?, ?, ?)",
		h.Name,
		h.IP,
		h.Domain,
		h.Role,
	)

	return err
}

func UpdateHost(db *sql.DB, h Host) error {

	_, err := db.Exec(
		`UPDATE hosts 
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
		"DELETE FROM hosts WHERE id = ?",
		id,
	)

	return err
}

func AddSubDomain(db *sql.DB, hostID int, subdomain string) error {

	_, err := db.Exec(
		"INSERT INTO subdomains(host_id, subdomain) VALUES (?, ?)",
		hostID,
		subdomain,
	)

	return err
}

func UpdateSubDomain(db *sql.DB, id int, subdomain string) error {

	_, err := db.Exec(
		"UPDATE subdomains SET subdomain = ? WHERE id = ?",
		subdomain,
		id,
	)

	return err
}

func DeleteSubDomain(db *sql.DB, id int) error {

	_, err := db.Exec(
		"DELETE FROM subdomains WHERE id = ?",
		id,
	)

	return err
}
