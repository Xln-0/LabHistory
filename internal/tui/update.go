package tui

import (
	"github.com/Xln-0/labhistory/internal/db"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		return m.updateKeys(msg)
	}

	return m, nil
}

func (m *Model) updateKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	// -------------------------
	// INPUT MODE (ADD / EDIT)
	// -------------------------
	if m.adding || m.editing {
		return m.updateEditor(msg)
	}

	if m.confirmDelete {
		return m.updateConfirmDelete(msg)
	}

	switch msg.String() {

	case "q":
		return m, tea.Quit

	case "tab":
		m.focus = (m.focus + 1) % 2

	case "a":
		m.startAdd()

	case "e":
		m.startEdit()

	case "d":
		m.startDelete()

	case "x":
		_ = m.ExportEnv()

	case "down":
		m.moveDown()

	case "up":
		m.moveUp()

	case "left":
		m.focus = 0

	case "right":
		m.focus = 1
	}

	return m, nil
}

// -------------------------
// START ADD
// -------------------------
func (m *Model) startAdd() {
	m.adding = true
	m.editing = false

	m.inputStep = 0

	m.input = ""
	m.cursorPos = 0

	m.formValues = make([]string, len(m.currentFields()))
}

// -------------------------
// START EDIT
// -------------------------
func (m *Model) startEdit() {
	if len(m.hosts) == 0 {
		return
	}

	m.editing = true
	m.adding = false

	m.editStep = 0

	m.tmp = m.hosts[m.hostSelected]

	m.input = m.tmp.Name
	m.cursorPos = len(m.input)
}

// -------------------------
// DELETE
// -------------------------
func (m *Model) startDelete() {
	if len(m.hosts) == 0 {
		return
	}

	m.confirmDelete = true
	m.deleteIndex = m.hostSelected
}

func (m *Model) updateConfirmDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	switch msg.String() {
	case "y":
		id := m.hosts[m.deleteIndex].ID
		_ = db.DeleteHost(m.db, id)
		m.loadHosts()

		if len(m.hosts) == 0 {
			m.hostSelected = 0
		} else if m.hostSelected >= len(m.hosts) {
			m.hostSelected = len(m.hosts) - 1
		}

	case "n", "esc":
	}

	m.confirmDelete = false

	return m, nil
}

// -------------------------
// NAVIGATION
// -------------------------
func (m *Model) moveUp() {
	if m.focus == 0 {
		if len(m.hosts) > 0 {
			m.hostSelected--
			if m.hostSelected < 0 {
				m.hostSelected = len(m.hosts) - 1
			}
		}
	} else {
		if len(m.users) > 0 {
			m.userSelected--
			if m.userSelected < 0 {
				m.userSelected = len(m.users) - 1
			}
		}
	}
}

func (m *Model) moveDown() {
	if m.focus == 0 {
		if len(m.hosts) > 0 {
			m.hostSelected = (m.hostSelected + 1) % len(m.hosts)
		}
	} else {
		if len(m.users) > 0 {
			m.userSelected = (m.userSelected + 1) % len(m.users)
		}
	}
}
