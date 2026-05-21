package tui

import (
	"database/sql"

	"github.com/Xln-0/labhistory/internal/db"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width  int
	height int

	mode Mode

	db *sql.DB

	hosts []db.Host
	users []db.User

	hostSelected int
	userSelected int

	focus int

	cursorPos int

	addForm  AddFormState
	editForm EditFormState

	deleteIndex int
}

func New(db *sql.DB) *Model {

	m := &Model{
		db:   db,
		mode: ModeNormal,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		loadHostsCmd(m.db),
		loadUsersCmd(m.db),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case hostsLoadedMsg:
		if msg.err != nil {
			// handle error (e.g., log it)
			return m, nil
		}

		m.hosts = msg.hosts
		return m, nil

	case usersLoadedMsg:
		if msg.err != nil {
			// handle error (e.g., log it)
			return m, nil
		}

		m.users = msg.users
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) View() string {

	header := m.renderHeader()

	body := m.renderBody()

	footer := m.renderFooter()

	return m.renderLayout(header, body, footer)
}
