package app

type PaletteCommand struct {
	Name   string
	Focus  Focus
	Action func(m *Model)
}

var PaletteCommands = []PaletteCommand{
	{
		Name:  "Add Host",
		Focus: FocusHosts,
		Action: func(m *Model) {
			m.startAddHost()
		},
	},
	{
		Name:  "Add Subdomain",
		Focus: FocusHosts,
		Action: func(m *Model) {
			m.startAddOption("subdomain")
		},
	},
	{
		Name:  "Delete Subdomain",
		Focus: FocusHosts,
		Action: func(m *Model) {
			m.startDeleteOption("subdomain")
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
