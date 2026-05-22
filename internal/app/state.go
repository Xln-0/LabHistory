package app

import (
	"strconv"

	"github.com/Xln-0/labhistory/internal/app/components"
)

// Focus
type Focus int

const (
	FocusHosts Focus = iota
	FocusUsers
)

// Mode
type Mode int

const (
	ModeNormal Mode = iota
	ModeAdd
	ModeEdit
	ModeDeleteConfirm
	ModePalette
	ModeAddOption
	ModeDeleteOption
	ModeSudoPrompt
)

type AddFormState struct {
	cursor int
	step   int

	values []string
	fields []string

	title  string
	target string // "host", "user", or specific field like "subdomain"
}

type EditSection struct {
	title string
	start int
	end   int

	list bool
}

type EditFormState struct {
	cursor int
	step   int

	fields []string
	values []string

	title  string
	target string

	target_id int

	sections []EditSection
}

type PaletteState struct {
	Cursor   int
	Title    string
	Commands components.List[PaletteCommand]
}

func (m Model) currentFields() []string {

	switch m.focus {

	// -------------------------
	// HOST FOCUS
	// -------------------------
	case FocusHosts:
		fields := []string{"Name", "IP", "Domain", "Role"}
		if m.mode != ModeAdd && m.hosts.Items[m.hosts.Selected].Subdomains != nil {
			for i := range m.hosts.Items[m.hosts.Selected].Subdomains {
				fields = append(fields, "Subdomain "+strconv.Itoa(i+1))
			}
		}
		return fields

	// -------------------------
	// USER FOCUS
	// -------------------------
	case FocusUsers:
		return []string{"Username", "Password", "Hash"}

	default:
		return []string{}
	}
}

func (m Model) currentValues() []string {

	switch m.focus {

	// -------------------------
	// HOST FOCUS
	// -------------------------
	case FocusHosts:
		host := m.hosts.Items[m.hosts.Selected]
		values := []string{host.Name, host.IP, host.Domain, host.Role}
		if host.Subdomains != nil {
			for _, s := range host.Subdomains {
				values = append(values, s.Subdomain)
			}
		}
		return values

	// -------------------------
	// USER FOCUS
	// -------------------------
	case FocusUsers:
		user := m.users.Items[m.users.Selected]
		return []string{user.Username, user.Password, user.Hash}

	default:
		return []string{}
	}
}

func (m Model) getHostSections() []EditSection {

	host := m.hosts.Items[m.hosts.Selected]

	sections := []EditSection{
		{title: "", start: 0, end: 4, list: false},
	}

	if len(host.Subdomains) > 0 {
		sections = append(sections, EditSection{
			title: "Subdomains",
			start: 4,
			end:   4 + len(host.Subdomains),
			list:  true,
		})
	}

	return sections
}

func (m Model) getUserSections() []EditSection {

	return []EditSection{
		{title: "", start: 0, end: 3, list: false},
	}
}
