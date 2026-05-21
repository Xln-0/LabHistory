package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeAdd
	ModeEdit
	ModeDeleteConfirm
)

type AddFormState struct {
	cursorPos int
	step      int
	input     string
	values    []string
}

type EditFormState struct {
	cursorPos int
	field     int
	fields    []string
	id        int
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	switch m.mode {

	case ModeNormal:
		return m.handleNormal(msg)

	case ModeAdd:
		return m.handleAdd(msg)

	case ModeEdit:
		return m.handleEdit(msg)

	case ModeDeleteConfirm:
		return m.handleDeleteConfirm(msg)
	}

	return m, nil
}
