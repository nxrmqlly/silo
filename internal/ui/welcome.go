package ui

import (
	"charm.land/lipgloss/v2"
	"github.com/nxrmqlly/silo/internal"
)

type Welcome struct {
	height  int
	width   int
	focused bool
}

func (w *Welcome) SetFocus(f bool) {
	w.focused = f
}

func (w *Welcome) SetSize(width, height int) {
	if width <= 2 || height <= 2 {
		return
	}
	w.width = width
	w.height = height
}

func (w *Welcome) View() string {
	if w.width == 0 || w.height == 0 {
		return ""
	}

	// border consumes 2 cols and 2 rows
	innerW := w.width - 2
	innerH := w.height - 2

	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	accentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("141")).Italic(true)
	purpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Bold(true)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		purpStyle.Render(internal.SiloAscii),
		purpStyle.Render(internal.WelcomeStr),
		"",
		accentStyle.Render("select a file to get started"),
		dimStyle.Render("n  new file · d  delete · tab  switch pane"),
	)

	centered := lipgloss.Place(
		innerW, innerH,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	borderColor := lipgloss.Color("238")
	if w.focused {
		borderColor = lipgloss.Color("139")
	}

	return lipgloss.NewStyle().
		Width(w.width).
		Height(w.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Render(centered)
}

func NewWelcome() *Welcome {
	return &Welcome{}
}
