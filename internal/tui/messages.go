package tui

import "github.com/Xln-0/labhistory/internal/db"

type hostsLoadedMsg struct {
	hosts []db.Host
	err   error
}

type usersLoadedMsg struct {
	users []db.User
	err   error
}
