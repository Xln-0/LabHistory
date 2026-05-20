package tui

import (
	"github.com/Xln-0/labhistory/internal/db"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateEditor(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.cursorPos > len(m.input) {
		m.cursorPos = len(m.input)
	}

	switch msg.String() {

	case "esc":
		m.adding = false
		m.editing = false
		m.input = ""
		m.cursorPos = 0
		m.inputStep = 0
		m.editStep = 0
		return m, nil
	case "enter":
		return m.handleEnter()

	case "backspace":
		return m.handleBackspace()

	case "left":
		if m.cursorPos > 0 {
			m.cursorPos--
		}

	case "right":
		if m.cursorPos < len(m.input) {
			m.cursorPos++
		}

	default:
		if m.cursorPos < 0 {
			m.cursorPos = 0
		}
		if m.cursorPos > len(m.input) {
			m.cursorPos = len(m.input)
		}
		if len(msg.String()) == 1 {
			m.input =
				m.input[:m.cursorPos] +
					msg.String() +
					m.input[m.cursorPos:]

			m.cursorPos++
		}
	}

	return m, nil
}

func (m *Model) handleEnter() (tea.Model, tea.Cmd) {

	if m.adding {
		return m.handleAddStep()
	}

	if m.editing {
		return m.handleEditStep()
	}

	return m, nil
}

func (m *Model) handleAddStep() (tea.Model, tea.Cmd) {
	m.formValues[m.inputStep] = m.input

	fields := m.currentFields()

	if m.inputStep == len(fields)-1 {

		if m.focus == 0 {

			h := db.Host{
				Name:   m.formValues[0],
				IP:     m.formValues[1],
				Domain: m.formValues[2],
				Role:   m.formValues[3],
			}

			_ = db.InsertHost(m.db, h)
			m.loadHosts()

		} else {

			u := db.User{
				Username: m.formValues[0],
				Password: m.formValues[1],
				Hash:     m.formValues[2],
			}

			_ = db.InsertUser(m.db, u)
			m.loadUsers()
		}

		m.adding = false
		m.inputStep = 0
		m.formValues = make([]string, len(fields))

		return m, nil
	}
	m.input = ""
	m.cursorPos = 0
	m.inputStep++

	return m, nil
}

func (m *Model) handleEditStep() (tea.Model, tea.Cmd) {
	switch m.editStep {
	case 0:
		m.tmp.Name = m.input
	case 1:
		m.tmp.IP = m.input
	case 2:
		m.tmp.Domain = m.input
	case 3:
		m.tmp.Role = m.input

		// FIN → update machine
		m.tmp.ID = m.hosts[m.hostSelected].ID
		_ = db.UpdateHost(m.db, m.tmp)
		m.loadHosts()
		m.editing = false
		m.input = ""
		m.cursorPos = 0
		m.editStep = 0
		return m, nil
	}

	m.editStep++
	m.input = m.currentFieldValue()
	m.cursorPos = len(m.input)
	return m, nil
}

func (m *Model) handleBackspace() (tea.Model, tea.Cmd) {
	if len(m.input) == 0 {
		return m, nil
	}

	if m.cursorPos > 0 {
		m.input =
			m.input[:m.cursorPos-1] +
				m.input[m.cursorPos:]

		m.cursorPos--
	}

	return m, nil
}

func (m Model) currentFields() []string {

	if m.focus == 0 {
		return []string{
			"Name",
			"IP",
			"Domain",
			"Role",
		}
	} else {
		return []string{
			"Username",
			"Password",
			"Hash",
		}
	}
}

func (m *Model) currentFieldValue() string {

	if len(m.hosts) == 0 {
		return ""
	}

	machine := m.hosts[m.hostSelected]

	if m.editStep == 0 {
		return machine.Name
	}
	if m.editStep == 1 {
		return machine.IP
	}
	if m.editStep == 2 {
		return machine.Domain
	}
	if m.editStep == 3 {
		return machine.Role
	}

	return ""
}
