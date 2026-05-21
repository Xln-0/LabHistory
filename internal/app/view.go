package app

import (
	"fmt"
	"strings"

	"github.com/Xln-0/labhistory/internal/app/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderLayout() string {

	header := m.renderHeader()
	body := m.renderBody()
	footer := m.renderFooter()

	if m.mode == ModePalette {
		body = m.renderPalette()
	}

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

func (m Model) renderBody() string {

	tabs := m.renderTabs()
	var middle string

	switch m.mode {

	case ModeAdd, ModeAddOption:
		middle = m.renderAddForm()

	case ModeEdit:
		middle = m.renderEditForm()

	case ModeDeleteConfirm:
		middle = m.renderConfirmDelete()

	case ModeDeleteOption:
		middle = m.renderDeleteOption()

	default:
		middle = m.renderDetails()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		tabs,
		middle,
	)
}

func (m Model) renderHosts() string {
	var style lipgloss.Style

	if m.focus == FocusHosts {
		style = styles.ActiveTextStyle
	} else {
		style = styles.InactiveTextStyle
	}

	return m.hosts.Render(m.width/2-2, style)
}

func (m Model) renderUsers() string {
	var style lipgloss.Style

	if m.focus == FocusUsers {
		style = styles.ActiveTextStyle
	} else {
		style = styles.InactiveTextStyle
	}

	return m.users.Render(m.width/2-2, style)
}

func (m Model) renderTabs() string {
	height := max(len(m.hosts.Items)+5, len(m.users.Items)+5)

	host := m.renderHosts()
	user := m.renderUsers()

	if m.mode != ModeNormal {
		host = styles.InactiveTab.Height(height).Width(m.width/2 - 2).Render(host)
		user = styles.InactiveTab.Height(height).Width(m.width/2 - 2 + m.width%2).Render(user)
	} else if m.focus == 0 {
		host = styles.ActiveTab.Height(height).Width(m.width/2 - 2).Render(host)
		user = styles.InactiveTab.Height(height).Width(m.width/2 - 2 + m.width%2).Render(user)
	} else {
		host = styles.InactiveTab.Height(height).Width(m.width/2 - 2).Render(host)
		user = styles.ActiveTab.Height(height).Width(m.width/2 - 2 + m.width%2).Render(user)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		host,
		user,
	)
}

func (m Model) renderDetails() string {

	content := "\n" + styles.Title.Render("DETAILS") + "\n───────\n\n"

	switch m.focus {

	// -------------------------
	// HOST FOCUS
	// -------------------------
	case FocusHosts:

		if len(m.hosts.Items) == 0 || m.hosts.Selected >= len(m.hosts.Items) {
			return styles.Border.Width(m.width - 2).Render("No host selected")
		}

		host := m.hosts.Items[m.hosts.Selected]

		content += "Name: " + host.Name + "\n"
		content += "IP: " + host.IP + "\n"
		content += "Domain: " + host.Domain + "\n"
		content += "Role: " + host.Role + "\n"

		if len(host.Subdomains) > 0 {
			content += "\nSubdomains:\n"

			for _, s := range host.Subdomains {
				content += " - " + s.Subdomain + "\n"
			}
		}

	// -------------------------
	// USER FOCUS
	// -------------------------
	case FocusUsers:

		if len(m.users.Items) == 0 || m.users.Selected >= len(m.users.Items) {
			return styles.Border.Width(m.width - 2).Render("No user selected")
		}

		user := m.users.Items[m.users.Selected]

		content += "Username: " + user.Username + "\n"
		content += "Password: " + user.Password + "\n"
		content += "Hash: " + user.Hash + "\n"
	}

	return styles.Border.Width(m.width - 2).Render(content)
}

func (m Model) renderAddForm() string {

	step := m.addForm.step

	if step >= len(m.addForm.fields) {
		step = len(m.addForm.fields) - 1
	}

	val := lipgloss.NewStyle().
		Render(renderWithCursor(m.addForm.values[step], m.addForm.cursor))

	out := "\n" + styles.Title.Render(m.addForm.title) + "\n"
	out += strings.Repeat("─", len(m.addForm.title)) + "\n\n"

	out += m.addForm.fields[step] + ": " + val + "\n\n"

	return styles.ActiveTab.Width(m.width - 2).Render(out)
}

func (m Model) renderEditForm() string {

	out := "\n" + styles.Title.Render(m.editForm.title) + "\n"
	out += strings.Repeat("─", len(m.editForm.title)) + "\n\n"

	for _, section := range m.editForm.sections {

		if section.list {
			out += "\n" + section.title + ":\n"
		}

		for i := section.start; i < section.end; i++ {

			val := m.editForm.values[i]

			if i == m.editForm.step {
				val = styles.ActiveTextStyle.
					Render(renderWithCursor(val, m.editForm.cursor))
			}

			if section.list {
				out += " - " + val + "\n"
			} else {
				out += fmt.Sprintf(
					"%s: %s\n",
					m.editForm.fields[i],
					val,
				)
			}
		}

	}

	return styles.ActiveTab.Width(m.width - 2).Render(out)
}

func (m Model) renderDeleteOption() string {

	out := m.option.Render(m.width-2, styles.Warning)

	return styles.ActiveTab.Width(m.width - 2).Render(out)
}

func (m Model) renderConfirmDelete() string {

	out := "\n"

	switch m.focus {

	case 0: // HOST
		h := m.hosts.Items[m.hosts.Selected]

		out += styles.Warning.Render("⚠ DELETE HOST ⚠") + "\n"
		out += strings.Repeat("─", len("  DELETE HOST  ")) + "\n\n"
		out += "Name: " + h.Name + "\n"
		out += "IP: " + h.IP + "\n\n"

	case 1: // USER
		u := m.users.Items[m.users.Selected]

		out += styles.Warning.Render("⚠ DELETE USER ⚠") + "\n"
		out += strings.Repeat("─", len("  DELETE USER  ")) + "\n\n"
		out += "Username: " + u.Username + "\n\n"
	}

	out += "Confirm delete? (y/n)\n"

	return styles.ActiveTab.Width(m.width - 2).Render(out)
}

func renderWithCursor(s string, cursor int) string {
	if cursor > len(s) {
		cursor = len(s)
	}

	return s[:cursor] + "|" + s[cursor:]
}

func (m *Model) renderHeader() string {

	ascii := `                                  
 __        _   _____ _     _               
|  |   ___| |_|  |  |_|___| |_ ___ ___ _ _ 
|  |__| .'| . |     | |_ -|  _| . |  _| | |
|_____|__,|___|__|__|_|___|_| |___|_| |_  |
                             by Xln-0 |___|
`

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9fef00")).
		Width(m.width).
		Padding(0, 0).
		Align(lipgloss.Center).
		Render(ascii)
}

func (m Model) renderFooter() string {

	var text string

	switch m.mode {
	case ModeDeleteConfirm:
		text = "y:confirm  n:cancel"

	case ModeAdd, ModeEdit:
		text = "enter:next/save  esc:cancel"

	case ModePalette:
		text = "↑/↓:navigate  enter:select  esc:cancel"

	default:
		text = "↑/↓:navigate  ←/→:panel  a:add  e:edit  d:delete  x:export  p:palette  q:quit"
	}

	return styles.FooterStyle.
		Width(m.width).
		Align(lipgloss.Center).
		Render(text)
}

func (m Model) renderPalette() string {
	var out string

	out += "\n" + styles.Title.Render("COMMAND PALETTE") + "\n"
	out += strings.Repeat("─", len("COMMAND PALETTE")) + "\n\n"

	for i, cmd := range m.palette.commands {

		cursor := " "
		if i == m.palette.cursor {
			cursor = ">"
		}

		out += cursor + " " + cmd.label + "\n"
	}

	return styles.ActiveTab.Width(m.width - 2).Render(out)
}
