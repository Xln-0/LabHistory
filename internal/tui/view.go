package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderLayout(header, body, footer string) string {
	spacer := lipgloss.NewStyle().
		Height(max(0, m.height-lipgloss.Height(header)-lipgloss.Height(body)-lipgloss.Height(footer))).
		Render("")

	screen := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
		spacer,
		footer,
	)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Render(screen)
}

func (m Model) renderHosts() string {

	out := "\nHOSTS\n─────\n\n"

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

	out := "\nUSERS\n─────\n\n"

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

func (m Model) renderHeader() string {
	height := max(len(m.hosts)+5, len(m.users)+5)

	host := m.renderHosts()
	user := m.renderUsers()

	if m.mode != ModeNormal {
		host = InactiveTab.Height(height).Width(m.width/2 - 2).Render(host)
		user = InactiveTab.Height(height).Width(m.width/2 - 2).Render(user)
	} else if m.focus == 0 {
		user = InactiveTab.Height(height).Width(m.width/2 - 2).Render(user)
		host = ActiveTab.Height(height).Width(m.width/2 - 2).Render(host)
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

func (m Model) renderBody() string {
	switch m.mode {

	case ModeAdd:
		return m.renderAddForm()

	case ModeEdit:
		return m.renderEditForm()

	case ModeDeleteConfirm:
		return m.renderConfirmDelete()

	default:
		return m.renderDetails()
	}
}

func (m Model) renderDetails() string {

	content := "\n" + Title.Render("DETAILS") + "\n───────\n\n"

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

	}

	// -------------------------
	// USER FOCUS
	// -------------------------
	if m.focus == 1 {

		if len(m.users) == 0 || m.userSelected >= len(m.users) {
			return Border.Width(m.width - 2).Render("No user selected")
		}

		user := m.users[m.userSelected]

		content += "Username: " + user.Username + "\n"
		content += "Password: " + user.Password + "\n"
		content += "Hash: " + user.Hash + "\n"
	}

	return Border.Width(m.width - 2).Render(content)
}

func (m Model) renderAddForm() string {
	cursor := "|"

	mode := "ADD HOST"
	fields := []string{"Host", "IP", "Domain", "Role"}
	if m.focus == 1 {
		mode = "ADD USER"
		fields = []string{"Username", "Password", "Hash"}
	}

	step := m.addForm.step

	if step >= len(fields) {
		step = len(fields) - 1
	}

	display := m.addForm.values[step]

	if m.cursorPos <= len(display) {
		display =
			display[:m.cursorPos] +
				cursor +
				display[m.cursorPos:]
	} else {
		display += cursor
	}

	out := "\n" + Title.Render(mode) + "\n"
	out += strings.Repeat("─", len(mode)) + "\n\n"

	out += fields[step] + ": " + display + "\n\n"

	return ActiveTab.Width(m.width - 2).Render(out)
}

func (m Model) renderEditForm() string {

	mode := "EDIT HOST"
	fields := []string{"Host", "IP", "Domain", "Role"}
	if m.focus == 1 {
		mode = "EDIT USER"
		fields = []string{"Username", "Password", "Hash"}
	}

	out := "\n" + Title.Render(mode) + "\n"
	out += strings.Repeat("─", len(mode)) + "\n\n"

	for i, label := range fields {

		val := m.editForm.fields[i]

		if i == m.editForm.field {
			val = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9fef00")).
				Render(renderWithCursor(val, m.cursorPos))
		}

		out += fmt.Sprintf("%s: %s\n", label, val)
	}

	return ActiveTab.Width(m.width - 2).Render(out)
}

func (m Model) renderConfirmDelete() string {

	out := "\n"

	switch m.focus {

	case 0: // HOST
		h := m.hosts[m.deleteIndex]

		out += Warning.Render("⚠ DELETE HOST ⚠") + "\n"
		out += strings.Repeat("─", len("  DELETE HOST  ")) + "\n\n"
		out += "Name: " + h.Name + "\n"
		out += "IP:   " + h.IP + "\n\n"

	case 1: // USER
		u := m.users[m.deleteIndex]

		out += Warning.Render("⚠ DELETE USER ⚠") + "\n"
		out += strings.Repeat("─", len("  DELETE USER  ")) + "\n\n"
		out += "Username: " + u.Username + "\n\n"
	}

	out += "Confirm delete? (y/n)\n"

	return ActiveTab.Width(m.width - 2).Render(out)
}

func renderWithCursor(s string, cursor int) string {
	if cursor > len(s) {
		cursor = len(s)
	}

	return s[:cursor] + "|" + s[cursor:]
}

func (m Model) renderFooter() string {

	var text string

	switch {
	case m.mode == ModeDeleteConfirm:
		text = "y:confirm  n:cancel"

	case m.mode == ModeAdd || m.mode == ModeEdit:
		text = "enter:next field/save  esc:cancel"

	default:
		text = "←/→:panel  a:add  e:edit  d:delete  x:export  q:quit"
	}

	return footerStyle.
		Width(m.width).
		Align(lipgloss.Center).
		Render(text)
}
