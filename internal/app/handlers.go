package app

import (
	"github.com/Xln-0/labhistory/internal/db"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	switch msg.String() {

	// -------------------------
	// QUIT
	// -------------------------
	case "q", "ctrl+c":
		return m, tea.Quit

	// -------------------------
	// SWITCH PANEL (hostS / userS)
	// -------------------------
	case "tab", "h", "l", "left", "right":
		m.focus = (m.focus + 1) % 2

	// -------------------------
	// NAVIGATION DOWN
	// -------------------------
	case "j", "down":
		if m.focus == 0 && len(m.hosts) > 0 {
			m.hostSelected = (m.hostSelected + 1) % len(m.hosts)
		}

		if m.focus == 1 && len(m.users) > 0 {
			m.userSelected = (m.userSelected + 1) % len(m.users)
		}

	// -------------------------
	// NAVIGATION UP
	// -------------------------
	case "k", "up":
		if m.focus == 0 && len(m.hosts) > 0 {
			m.hostSelected--
			if m.hostSelected < 0 {
				m.hostSelected = len(m.hosts) - 1
			}
		}

		if m.focus == 1 && len(m.users) > 0 {
			m.userSelected--
			if m.userSelected < 0 {
				m.userSelected = len(m.users) - 1
			}
		}

	// -------------------------
	// CREATE
	// -------------------------
	case "a":
		m.mode = ModeAdd
		m.addForm.input = ""
		m.addForm.step = 0
		m.addForm.cursorPos = 0
		m.addForm.values = make([]string, len(m.currentFields()))

	// -------------------------
	// EDIT
	// -------------------------
	case "e":
		m.mode = ModeEdit

		if m.focus == 0 && len(m.hosts) > 0 {
			h := m.hosts[m.hostSelected]

			m.editForm.id = h.ID

			m.editForm.fields = []string{
				h.Name,
				h.IP,
				h.Domain,
				h.Role,
			}
		}

		if m.focus == 1 && len(m.users) > 0 {
			u := m.users[m.userSelected]

			m.editForm.id = u.ID

			m.editForm.fields = []string{
				u.Username,
				u.Password,
				u.Hash,
			}
		}

		m.editForm.field = 0
		m.cursorPos = len(m.editForm.fields[m.editForm.field])

	// -------------------------
	// DELETE
	// -------------------------
	case "d":
		m.mode = ModeDeleteConfirm

		if m.focus == 0 && len(m.hosts) > 0 {
			m.deleteIndex = m.hostSelected
		}
		if m.focus == 1 && len(m.users) > 0 {
			m.deleteIndex = m.userSelected
		}

	// -------------------------
	// EXPORT
	// -------------------------
	case "x":
		_ = m.ExportEnv()
	}

	return m, nil
}

func (m *Model) handleAdd(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	switch msg.String() {

	// -------------------------
	// CANCEL
	// -------------------------
	case "esc":
		m.mode = ModeNormal

	// -------------------------
	// NEXT FIELD
	// -------------------------
	case "enter":
		if m.addForm.step == len(m.currentFields())-1 {
			m.mode = ModeNormal

			if m.focus == 0 {
				// Create Host
				h := db.Host{
					Name:   m.addForm.values[0],
					IP:     m.addForm.values[1],
					Domain: m.addForm.values[2],
					Role:   m.addForm.values[3],
				}
				_ = db.CreateHost(m.db, h)
				return m, loadHostsCmd(m.db)
			}
			if m.focus == 1 {
				// Create User
				u := db.User{
					Username: m.addForm.values[0],
					Password: m.addForm.values[1],
					Hash:     m.addForm.values[2],
				}
				_ = db.CreateUser(m.db, u)
				return m, loadUsersCmd(m.db)
			}
		}
		m.addForm.step++
		m.addForm.input = ""
		m.cursorPos = 0

	// -------------------------
	// BACKSPACE
	// -------------------------
	case "backspace":
		val := m.addForm.values[m.addForm.step]
		if m.cursorPos > 0 && len(val) > 0 {
			m.addForm.values[m.addForm.step] =
				val[:m.cursorPos-1] +
					val[m.cursorPos:]

			m.cursorPos--
		}

	// -------------------------
	// LEFT
	// -------------------------
	case "left":
		if m.cursorPos > 0 {
			m.cursorPos--
		}

	// -------------------------
	// RIGHT
	// -------------------------
	case "right":
		val := m.addForm.values[m.addForm.step]
		if m.cursorPos < len(val) {
			m.cursorPos++
		}

	// -------------------------
	// UPDATE INPUT
	// -------------------------
	default:
		if msg.Type == tea.KeyRunes {
			val := m.addForm.values[m.addForm.step]

			m.addForm.values[m.addForm.step] =
				val[:m.cursorPos] +
					msg.String() +
					val[m.cursorPos:]

			m.cursorPos++
		}
	}

	return m, nil
}

func (m *Model) handleEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	switch msg.String() {

	// -------------------------
	// CANCEL
	// -------------------------
	case "esc":
		m.mode = ModeNormal

	// -------------------------
	// ENTER
	// -------------------------
	case "enter":

		if m.editForm.field == len(m.editForm.fields)-1 {
			m.mode = ModeNormal

			if m.focus == 0 {

				// Update Host
				h := db.Host{
					ID:     m.editForm.id,
					Name:   m.editForm.fields[0],
					IP:     m.editForm.fields[1],
					Domain: m.editForm.fields[2],
					Role:   m.editForm.fields[3],
				}
				_ = db.UpdateHost(m.db, h)
				return m, loadHostsCmd(m.db)
			}

			if m.focus == 1 {

				// Update User
				u := db.User{
					ID:       m.editForm.id,
					Username: m.editForm.fields[0],
					Password: m.editForm.fields[1],
					Hash:     m.editForm.fields[2],
				}
				_ = db.UpdateUser(m.db, u)
				return m, loadUsersCmd(m.db)
			}
		}

		m.editForm.field++
		m.cursorPos = len(m.editForm.fields[m.editForm.field])

	// -------------------------
	// BACKSPACE
	// -------------------------
	case "backspace":
		val := m.editForm.fields[m.editForm.field]

		if m.cursorPos > 0 && len(val) > 0 {
			m.editForm.fields[m.editForm.field] =
				val[:m.cursorPos-1] +
					val[m.cursorPos:]

			m.cursorPos--
		}

	// -------------------------
	// LEFT
	// -------------------------
	case "left":
		if m.cursorPos > 0 {
			m.cursorPos--
		}

	// -------------------------
	// RIGHT
	// -------------------------
	case "right":
		val := m.editForm.fields[m.editForm.field]
		if m.cursorPos < len(val) {
			m.cursorPos++
		}

	// -------------------------
	// UP
	// -------------------------
	case "up":
		if m.editForm.field > 0 {
			m.editForm.field--
			m.cursorPos = len(m.editForm.fields[m.editForm.field])
		}

	// -------------------------
	// DOWN
	// -------------------------
	case "down":
		if m.editForm.field < len(m.editForm.fields)-1 {
			m.editForm.field++
			m.cursorPos = len(m.editForm.fields[m.editForm.field])
		}

	// -------------------------
	// UPDATE INPUT
	// -------------------------
	default:
		if msg.Type == tea.KeyRunes {
			val := m.editForm.fields[m.editForm.field]

			val =
				val[:m.cursorPos] +
					msg.String() +
					val[m.cursorPos:]

			m.editForm.fields[m.editForm.field] = val
			m.cursorPos++
		}
	}

	return m, nil
}

func (m *Model) handleDeleteConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	switch msg.String() {

	case "y":

		switch m.focus {

		case 0: // HOSTS
			id := m.hosts[m.deleteIndex].ID
			_ = db.DeleteHost(m.db, id)

			if len(m.hosts) == 0 {
				m.hostSelected = 0
			} else if m.hostSelected >= len(m.hosts)-1 {
				m.hostSelected = len(m.hosts) - 2
			}

			return m, loadHostsCmd(m.db)

		case 1: // USERS
			id := m.users[m.deleteIndex].ID
			_ = db.DeleteUser(m.db, id)

			if len(m.users) == 0 {
				m.userSelected = 0
			} else if m.userSelected >= len(m.users)-1 {
				m.userSelected = len(m.users) - 2
			}

			return m, loadUsersCmd(m.db)

		}

	}

	m.mode = ModeNormal

	return m, nil
}
