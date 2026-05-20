package tui

import "github.com/charmbracelet/lipgloss"

func (m Model) renderHosts() string {

	out := "HOSTS\n─────\n\n"

	if len(m.hosts) == 0 {
		out += "No Hosts"
		return Border.Width(m.width/2 - 2).Render(out)
	}

	for i, h := range m.hosts {

		cursor := " "
		if i == m.hostSelected && m.focus == 0 {
			cursor = ">"
		}

		out += cursor + " " + h.Name + "\n"
	}

	return out
}

func (m Model) renderUsers() string {

	out := "USERS\n─────\n\n"

	if len(m.users) == 0 {
		out += "No Users"
		return Border.Width(m.width/2 - 2).Render(out)
	}

	for i, u := range m.users {

		cursor := " "
		if i == m.userSelected && m.focus == 1 {
			cursor = ">"
		}

		out += cursor + " " + u.Username + "\n"
	}

	return out
}

func (m Model) renderTop(height int) string {
	host := m.renderHosts()
	user := m.renderUsers()

	if m.focus == 0 {
		host = ActiveTab.Height(height).Width(m.width/2 - 2).Render(host)
		user = InactiveTab.Height(height).Width(m.width/2 - 2).Render(user)
	} else {
		host = InactiveTab.Height(height).Width(m.width/2 - 2).Render(host)
		user = ActiveTab.Height(height).Width(m.width/2 - 2).Render(user)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		host,
		user,
	)
}

func (m Model) renderDetails() string {

	content := Title.Render("DETAILS") + "\n───────\n\n"

	// -------------------------
	// HOST FOCUS
	// -------------------------
	if m.focus == 0 {

		if len(m.hosts) == 0 || m.hostSelected >= len(m.hosts) {
			return Border.Width(m.width - 2).Render("No host selected")
		}

		host := m.hosts[m.hostSelected]

		content += "Host: " + host.Name + "\n"
		content += "IP: " + host.IP + "\n"
		content += "Domain: " + host.Domain + "\n"
		content += "Role: " + host.Role + "\n"

		return Border.Width(m.width - 2).Render(content)
	}

	// -------------------------
	// USER FOCUS
	// -------------------------
	if m.focus == 1 {

		if len(m.users) == 0 || m.userSelected >= len(m.users) {
			return Border.Width(m.width - 2).Render("No user selected")
		}

		user := m.users[m.userSelected]

		content += "User: " + user.Username + "\n"
		content += "Password: " + user.Password + "\n"
		content += "Hash: " + user.Hash + "\n"

		return Border.Width(m.width - 2).Render(content)
	}

	return ""
}

func (m Model) renderInput() string {
	cursor := "|"

	mode := "HOST"
	fields := []string{"Name", "IP", "Domain", "Role"}
	if m.focus == 1 {
		mode = "USER"
		fields = []string{"Username", "Password", "Hash"}
	}

	action := "ADD"
	if m.editing {
		action = "EDIT"
	}

	step := m.inputStep
	if m.editing {
		step = m.editStep
	}

	if step >= len(fields) {
		step = len(fields) - 1
	}

	display := m.input

	if m.cursorPos <= len(display) {
		display =
			display[:m.cursorPos] +
				cursor +
				display[m.cursorPos:]
	} else {
		display += cursor
	}

	out := action + " " + mode + "\n"
	out += "────────────\n\n"

	out += fields[step] + ": " + display + "\n\n"

	return Border.Width(40).Render(out)
}

func (m Model) renderConfirmDelete() string {

	h := m.hosts[m.deleteIndex]

	out := "⚠ DELETE HOST ⚠\n\n"
	out += "Name: " + h.Name + "\n"
	out += "IP:   " + h.IP + "\n\n"
	out += "Confirm delete? (y/n)\n"

	return Border.Width(40).Render(out)
}

func (m Model) renderFooter() string {

	var text string

	switch {
	case m.confirmDelete:
		text = "y:confirm  n:cancel"

	case m.adding || m.editing:
		text = "enter:next field  esc:cancel"

	default:
		text = "←/→:panel  a:add  e:edit  d:delete  x:export  q:quit"
	}

	return footerStyle.
		Width(m.width).
		Align(lipgloss.Center).
		Render(text)
}
