package app

import (
	"database/sql"

	"github.com/Xln-0/labhistory/internal/app/components"
	"github.com/Xln-0/labhistory/internal/db"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width  int
	height int

	focus  Focus
	mode   Mode
	status string

	db *sql.DB

	hosts  components.List[db.Host]
	users  components.List[db.User]
	option components.List[components.ListItem]

	addForm  AddFormState
	editForm EditFormState

	sudoInput textinput.Model
	sudoCmd   func(password string) tea.Cmd

	Palette components.List[PaletteCommand]

	errMessage string
}

func New(db *sql.DB) *Model {

	m := &Model{
		db:    db,
		mode:  ModeNormal,
		focus: FocusHosts,
	}

	m.hosts.Title = "hosts"
	m.users.Title = "users"

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
			m.errMessage = msg.err.Error()
			return m, nil
		}

		m.hosts.Items = msg.hosts
		return m, nil

	case usersLoadedMsg:
		if msg.err != nil {
			m.errMessage = msg.err.Error()
			return m, nil
		}

		m.users.Items = msg.users
		return m, nil

	case syncFinishedMsg:
		if msg.err != nil {
			m.errMessage = msg.err.Error()
		} else {
			m.errMessage = "Sync OK"
		}

		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) View() string {
	return m.renderLayout()
}
