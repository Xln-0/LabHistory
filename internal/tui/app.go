package tui

import (
	"database/sql"

	"github.com/Xln-0/labhistory/internal/db"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width  int
	height int

	focus     int // 0 = hosts, 1 = details
	cursorPos int

	// data
	hosts        []db.Host
	hostSelected int
	users        []db.User
	userSelected int

	// add mode
	adding     bool
	inputStep  int
	input      string
	formValues []string
	tmp        db.Host

	// edit mode
	editing  bool
	editStep int

	// delete
	confirmDelete bool
	deleteIndex   int

	// db
	db *sql.DB
}

func New(db *sql.DB) *Model {

	m := &Model{
		db: db,
	}

	m.loadHosts()
	m.loadUsers()

	return m
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) loadHosts() error {

	list, err := db.LoadHosts(m.db)
	if err != nil {
		return err
	}

	m.hosts = list
	return nil
}

func (m *Model) loadUsers() error {

	list, err := db.LoadUsers(m.db)
	if err != nil {
		return err
	}

	m.users = list
	return nil
}

func (m *Model) View() string {

	if m.confirmDelete {
		return m.renderConfirmDelete()
	}

	var body string

	if m.adding || m.editing {
		body = m.renderInput()
	} else {
		panelHeight := max(len(m.hosts)+4, len(m.users)+4)
		top := m.renderTop(panelHeight)
		middle := m.renderDetails()

		body = lipgloss.JoinVertical(
			lipgloss.Left,
			top,
			middle,
		)
	}

	footer := m.renderFooter()

	spacer := lipgloss.NewStyle().
		Height(max(0, m.height-lipgloss.Height(body)-lipgloss.Height(footer))).
		Render("")

	screen := lipgloss.JoinVertical(
		lipgloss.Left,
		body,
		spacer,
		footer,
	)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Render(screen)
}
