package app

import (
	"strings"

	"github.com/Xln-0/labhistory/internal/app/components"
	"github.com/Xln-0/labhistory/internal/db"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	switch m.mode {

	case ModeNormal:
		return m.handleNormal(msg)

	case ModeAdd, ModeAddOption:
		return m.handleAdd(msg)

	case ModeEdit:
		return m.handleEdit(msg)

	case ModeDeleteConfirm:
		return m.handleDeleteConfirm(msg)

	case ModePalette:
		return m.handlePalette(msg)

	case ModeDeleteOption:
		return m.handleDeleteOption(msg)

	case ModeSudoPrompt:
		return m.handleSudoPrompt(msg)

	}

	return m, nil
}

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
		if m.focus == 0 && len(m.hosts.Items) > 0 {
			m.hosts.Selected = (m.hosts.Selected + 1) % len(m.hosts.Items)
		}

		if m.focus == 1 && len(m.users.Items) > 0 {
			m.users.Selected = (m.users.Selected + 1) % len(m.users.Items)
		}

	// -------------------------
	// NAVIGATION UP
	// -------------------------
	case "k", "up":
		if m.focus == 0 && len(m.hosts.Items) > 0 {
			m.hosts.Selected--
			if m.hosts.Selected < 0 {
				m.hosts.Selected = len(m.hosts.Items) - 1
			}
		}

		if m.focus == 1 && len(m.users.Items) > 0 {
			m.users.Selected--
			if m.users.Selected < 0 {
				m.users.Selected = len(m.users.Items) - 1
			}
		}

	// -------------------------
	// CREATE
	// -------------------------
	case "a":
		if m.focus == 0 {
			m.startAddHost()
		} else {
			m.startAddUser()
		}

	// -------------------------
	// EDIT
	// -------------------------
	case "e":
		if m.focus == 0 && len(m.hosts.Items) > 0 {
			m.startEditHost()
		}

		if m.focus == 1 && len(m.users.Items) > 0 {
			m.startEditUser()
		}

		m.editForm.cursor = len(m.editForm.values[m.editForm.step])

	// -------------------------
	// DELETE
	// -------------------------
	case "d":
		if (m.focus == 0 && len(m.hosts.Items) > 0) || (m.focus == 1 && len(m.users.Items) > 0) {
			m.mode = ModeDeleteConfirm
		}

	// -------------------------
	// EXPORT
	// -------------------------
	case "x":
		_ = m.ExportEnv()

	// -------------------------
	// Palette
	// -------------------------
	case "p", "ctrl+p":
		m.mode = ModePalette
		m.Palette.Selected = 0
		m.Palette.Items = m.getFocusedCommands()
		if m.focus == FocusHosts {
			m.Palette.Title = "HOST COMMANDS PALETTE"
		} else {
			m.Palette.Title = "USER COMMANDS PALETTE"
		}
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
		if m.addForm.step == len(m.addForm.fields)-1 {

			m.mode = ModeNormal

			switch m.addForm.target {

			case "host":
				// Create Host
				h := db.Host{
					Name:   m.addForm.values[0],
					IP:     m.addForm.values[1],
					Domain: m.addForm.values[2],
					Role:   m.addForm.values[3],
				}
				_ = db.CreateHost(m.db, h)
				return m, loadHostsCmd(m.db)

			case "user":
				// Create User
				u := db.User{
					Username: m.addForm.values[0],
					Password: m.addForm.values[1],
					Hash:     m.addForm.values[2],
				}
				_ = db.CreateUser(m.db, u)
				return m, loadUsersCmd(m.db)

			case "subdomain":
				host := m.hosts.Items[m.hosts.Selected]

				_ = db.AddSubDomain(m.db, host.ID, m.addForm.values[0])
				return m, loadHostsCmd(m.db)
			}
		}
		m.addForm.step++
		m.addForm.cursor = 0

	// -------------------------
	// BACKSPACE
	// -------------------------
	case "backspace":
		val := m.addForm.values[m.addForm.step]
		if m.addForm.cursor > 0 && len(val) > 0 {
			m.addForm.values[m.addForm.step] =
				val[:m.addForm.cursor-1] +
					val[m.addForm.cursor:]

			m.addForm.cursor--
		}

	// -------------------------
	// LEFT
	// -------------------------
	case "left":
		if m.addForm.cursor > 0 {
			m.addForm.cursor--
		}

	// -------------------------
	// RIGHT
	// -------------------------
	case "right":
		val := m.addForm.values[m.addForm.step]
		if m.addForm.cursor < len(val) {
			m.addForm.cursor++
		}

	// -------------------------
	// UPDATE INPUT
	// -------------------------
	default:
		if msg.Type == tea.KeyRunes {
			val := m.addForm.values[m.addForm.step]

			m.addForm.values[m.addForm.step] =
				val[:m.addForm.cursor] +
					msg.String() +
					val[m.addForm.cursor:]

			m.addForm.cursor++
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

		m.mode = ModeNormal

		if m.focus == 0 {

			// Update Host
			host := m.hosts.Items[m.hosts.Selected]
			newHost := db.Host{
				ID:     m.editForm.target_id,
				Name:   m.editForm.values[0],
				IP:     m.editForm.values[1],
				Domain: m.editForm.values[2],
				Role:   m.editForm.values[3],
			}
			if host.Name != newHost.Name ||
				host.IP != newHost.IP ||
				host.Domain != newHost.Domain ||
				host.Role != newHost.Role {
				_ = db.UpdateHost(m.db, newHost)
			}

			if len(m.editForm.sections) > 1 {
				// Update Options
				for _, section := range m.editForm.sections {
					switch section.title {
					case "Subdomains":
						host := m.hosts.Items[m.hosts.Selected]
						for i := range host.Subdomains {
							s := host.Subdomains[i]
							newVal := m.editForm.values[4+i]
							if s.Subdomain != newVal {
								_ = db.UpdateSubDomain(m.db, s.ID, newVal)
							}
						}
					}
				}
			}

			return m, loadHostsCmd(m.db)
		}

		if m.focus == 1 {

			// Update User
			u := db.User{
				ID:       m.editForm.target_id,
				Username: m.editForm.values[0],
				Password: m.editForm.values[1],
				Hash:     m.editForm.values[2],
			}
			_ = db.UpdateUser(m.db, u)
			return m, loadUsersCmd(m.db)
		}

	// -------------------------
	// BACKSPACE
	// -------------------------
	case "backspace":
		val := m.editForm.values[m.editForm.step]

		if m.editForm.cursor > 0 && len(val) > 0 {
			m.editForm.values[m.editForm.step] =
				val[:m.editForm.cursor-1] +
					val[m.editForm.cursor:]

			m.editForm.cursor--
		}

	// -------------------------
	// LEFT
	// -------------------------
	case "left":
		if m.editForm.cursor > 0 {
			m.editForm.cursor--
		}

	// -------------------------
	// RIGHT
	// -------------------------
	case "right":
		val := m.editForm.values[m.editForm.step]
		if m.editForm.cursor < len(val) {
			m.editForm.cursor++
		}

	// -------------------------
	// UP
	// -------------------------
	case "up":
		if m.editForm.step > 0 {
			m.editForm.step--
			m.editForm.cursor = len(m.editForm.values[m.editForm.step])
		}

	// -------------------------
	// DOWN
	// -------------------------
	case "down":
		if m.editForm.step < len(m.editForm.values)-1 {
			m.editForm.step++
			m.editForm.cursor = len(m.editForm.values[m.editForm.step])
		}

	// -------------------------
	// UPDATE INPUT
	// -------------------------
	default:
		if msg.Type == tea.KeyRunes {
			val := m.editForm.values[m.editForm.step]

			m.editForm.values[m.editForm.step] =
				val[:m.editForm.cursor] +
					msg.String() +
					val[m.editForm.cursor:]

			m.editForm.cursor++
		}
	}

	return m, nil
}

func (m *Model) handleDeleteConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	m.mode = ModeNormal

	switch msg.String() {

	case "y":

		switch m.focus {

		case FocusHosts: // HOSTS
			id := m.hosts.Items[m.hosts.Selected].ID
			_ = db.DeleteHost(m.db, id)

			if len(m.hosts.Items) == 0 {
				m.hosts.Selected = 0
			} else if m.hosts.Selected >= len(m.hosts.Items)-1 {
				m.hosts.Selected = len(m.hosts.Items) - 2
			}

			return m, loadHostsCmd(m.db)

		case FocusUsers: // USERS
			id := m.users.Items[m.users.Selected].ID
			_ = db.DeleteUser(m.db, id)

			if len(m.users.Items) == 0 {
				m.users.Selected = 0
			} else if m.users.Selected >= len(m.users.Items)-1 {
				m.users.Selected = len(m.users.Items) - 2
			}

			return m, loadUsersCmd(m.db)

		}

	}

	return m, nil
}

func (m *Model) handleDeleteOption(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	switch msg.String() {

	// -------------------------
	// CANCEL
	// -------------------------
	case "esc":
		m.mode = ModeNormal

	// -------------------------
	// NAVIGATION DOWN
	// -------------------------
	case "j", "down":
		if len(m.option.Items) > 0 {
			m.option.Selected = (m.option.Selected + 1) % len(m.option.Items)
		}

	// -------------------------
	// NAVIGATION UP
	// -------------------------
	case "k", "up":
		if len(m.option.Items) > 0 {
			m.option.Selected--
			if m.option.Selected < 0 {
				m.option.Selected = len(m.option.Items) - 1
			}
		}

	// -------------------------
	// SELECT
	// -------------------------
	case "enter":
		m.mode = ModeNormal

		id := m.option.Items[m.option.Selected].ID
		_ = db.DeleteSubDomain(m.db, id)

		return m, loadHostsCmd(m.db)
	}

	return m, nil
}

func (m *Model) handlePalette(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	switch msg.String() {

	// -------------------------
	// CANCEL
	// -------------------------
	case "esc", "p":
		m.mode = ModeNormal

	// -------------------------
	// NAVIGATION DOWN
	// -------------------------
	case "j", "down":
		if len(m.Palette.Items) > 0 {
			m.Palette.Selected = (m.Palette.Selected + 1) % len(m.Palette.Items)
		}

	// -------------------------
	// NAVIGATION UP
	// -------------------------
	case "k", "up":
		if len(m.Palette.Items) > 0 {
			m.Palette.Selected--
			if m.Palette.Selected < 0 {
				m.Palette.Selected = len(m.Palette.Items) - 1
			}
		}

	// -------------------------
	// SELECT
	// -------------------------
	case "enter":
		m.Palette.Items[m.Palette.Selected].Action(m)
	}

	return m, nil
}

func (m *Model) handleSudoPrompt(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	switch msg.String() {
	case "esc":
		m.mode = ModeNormal
		return m, nil

	case "enter":
		password := m.sudoInput.Value()

		m.mode = ModeNormal

		return m, m.sudoCmd(password)
	}

	var cmd tea.Cmd
	m.sudoInput, cmd = m.sudoInput.Update(msg)

	return m, cmd
}

func (m *Model) startAddHost() {
	m.mode = ModeAdd
	m.focus = 0
	m.addForm = AddFormState{
		cursor: 0,
		step:   0,
		values: make([]string, len(m.currentFields())),
		fields: m.currentFields(),
		title:  "ADD HOST",
		target: "host",
	}
}

func (m *Model) startAddUser() {
	m.mode = ModeAdd
	m.focus = 1
	m.addForm = AddFormState{
		cursor: 0,
		step:   0,
		values: make([]string, len(m.currentFields())),
		fields: m.currentFields(),
		title:  "ADD USER",
		target: "user",
	}
}

func (m *Model) startEditHost() {
	m.mode = ModeEdit
	m.editForm = EditFormState{
		cursor:    0,
		step:      0,
		values:    m.currentValues(),
		fields:    m.currentFields(),
		title:     "EDIT HOST",
		target:    "host",
		target_id: m.hosts.Items[m.hosts.Selected].ID,
		sections:  m.getHostSections(),
	}
}

func (m *Model) startEditUser() {
	m.mode = ModeEdit
	m.editForm = EditFormState{
		cursor:    0,
		step:      0,
		values:    m.currentValues(),
		fields:    m.currentFields(),
		title:     "EDIT USER",
		target:    "user",
		target_id: m.users.Items[m.users.Selected].ID,
		sections:  m.getUserSections(),
	}
}

func (m *Model) startAddOption(option string) {
	m.mode = ModeAddOption
	m.addForm = AddFormState{
		cursor: 0,
		step:   0,
		values: []string{""},
		fields: []string{CapitalizeFirst(option)},
		title:  "ADD " + strings.ToUpper(option),
		target: option,
	}
}

func (m *Model) startDeleteOption(option string) {

	if m.focus != FocusHosts {
		return
	}

	host := m.hosts.Items[m.hosts.Selected]

	switch option {
	case "subdomain":
		if len(host.Subdomains) == 0 {
			return
		}

		m.option = components.List[components.ListItem]{
			Items:    subdomainsToList(host.Subdomains),
			Selected: 0,
			Title:    option,
			Prefix:   "",
		}
	}

	m.mode = ModeDeleteOption
}

func (m *Model) initSudoPrompt() {
	ti := textinput.New()

	ti.Placeholder = "sudo password"
	ti.Focus()

	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'

	m.sudoInput = ti
}
