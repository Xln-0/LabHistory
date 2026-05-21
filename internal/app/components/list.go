package components

import (
	"strings"

	"github.com/Xln-0/labhistory/internal/app/styles"
	"github.com/charmbracelet/lipgloss"
)

type ListItem struct {
	ID   int
	Name string
}

func (l ListItem) Label() string {
	return l.Name
}

type Labeler interface {
	Label() string
}

type List[T Labeler] struct {
	Items    []T
	Selected int

	Title  string
	Prefix string
}

func (l List[T]) Render(width int, style lipgloss.Style) string {

	out := "\n" + styles.Title.Render(strings.ToUpper(l.Title)) + "\n"
	out += strings.Repeat("─", len(strings.ToUpper(l.Title))) + "\n\n"

	for i, item := range l.Items {

		line := l.Prefix + item.Label()

		if i == l.Selected {
			line = style.Render("> " + line)
		} else {
			line = "  " + line
		}

		out += line + "\n"
	}

	return lipgloss.NewStyle().
		Width(width).
		Render(out)
}
