package app

import (
	"database/sql"

	"github.com/Xln-0/labhistory/internal/db"
	tea "github.com/charmbracelet/bubbletea"
)

func loadHostsCmd(database *sql.DB) tea.Cmd {

	return func() tea.Msg {

		hosts, err := db.LoadHosts(database)

		return hostsLoadedMsg{
			hosts: hosts,
			err:   err,
		}
	}
}

func loadUsersCmd(database *sql.DB) tea.Cmd {

	return func() tea.Msg {

		users, err := db.LoadUsers(database)

		return usersLoadedMsg{
			users: users,
			err:   err,
		}
	}
}
