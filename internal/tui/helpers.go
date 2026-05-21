package tui

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
