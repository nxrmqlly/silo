package ui

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/glamour"
)

type Preview struct {
	viewport viewport.Model
	renderer *glamour.TermRenderer

	content string

	width   int
	height  int
	focused bool
}

func (p *Preview) SetFocus(f bool) {
	p.focused = f
}

func (p *Preview) SetContent(content string) {
	p.content = content
	if p.renderer != nil && p.content != "" {
		rendered, _ := p.renderer.Render(p.content)
		p.viewport.SetContent(rendered)
	}
}

func (p *Preview) SetSize(w, h int) {
	if w <= 2 || h <= 2 {
		return
	}

	p.width = w
	p.height = h

	p.viewport.SetWidth(w - 2)
	p.viewport.SetHeight(h - 2)

	wrapWidth := w - 4 //for padding

	if wrapWidth < 20 {
		wrapWidth = 20
	}

	rd, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(wrapWidth),
	)

	p.renderer = rd

	if p.content != "" {
		rendered, _ := p.renderer.Render(p.content)
		p.viewport.SetContent(rendered)
	}

}

func (p *Preview) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	return cmd
}

func (p *Preview) View() string {
	if p.width == 0 || p.height == 0 {
		return ""
	}

	style := lipgloss.NewStyle().
		Width(p.width).
		Height(p.height).
		Padding(0, 0)

	if p.focused {
		style = style.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("139"))
	} else {
		style = style.Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238"))
	}

	return style.Render(p.viewport.View())

}

func NewPreview() *Preview {
	vp := viewport.New(
		viewport.WithHeight(0),
		viewport.WithWidth(0),
	)

	rd, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80), // default
	)

	return &Preview{
		viewport: vp,
		renderer: rd,
		content:  "",
	}
}
