package app

import tea "github.com/charmbracelet/bubbletea"

type PaletteCommand struct {
	Name   string
	Focus  Focus
	Action func(m *Model) tea.Cmd
}

var PaletteCommands = []PaletteCommand{
	{
		Name:  "Add Host",
		Focus: FocusHosts,
		Action: func(m *Model) tea.Cmd {
			m.startAddHost()
			return nil
		},
	},
	{
		Name:  "Add User",
		Focus: FocusUsers,
		Action: func(m *Model) tea.Cmd {
			m.startAddUser()
			return nil
		},
	},
	{
		Name:  "Add Subdomain",
		Focus: FocusHosts,
		Action: func(m *Model) tea.Cmd {
			m.startAddOption("subdomain")
			return nil
		},
	},
	{
		Name:  "Delete Subdomain",
		Focus: FocusHosts,
		Action: func(m *Model) tea.Cmd {
			m.startDeleteOption("subdomain")
			return nil
		},
	},
	{
		Name:  "Sync /etc/hosts",
		Focus: FocusHosts,
		Action: func(m *Model) tea.Cmd {

			host := m.hosts.Items[m.hosts.Selected]

			line := buildHostsLine(
				host.IP,
				host.Domain,
				subdomainsToList(host.Subdomains),
			)

			m.initSudoPrompt()

			m.sudoCmd = func(password string) tea.Cmd {
				return syncEtcHostsCmd(line, password)
			}

			m.mode = ModeSudoPrompt

			return nil
		},
	},
}

func (cmd PaletteCommand) Label() string {
	return cmd.Name
}

func (m *Model) getFocusedCommands() []PaletteCommand {

	out := make([]PaletteCommand, 0)

	for _, c := range PaletteCommands {
		if c.Focus == m.focus {
			out = append(out, c)
		}
	}

	return out
}
