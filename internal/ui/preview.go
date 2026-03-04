package ui

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/glamour"
)

type Preview struct {
	viewport      viewport.Model
	content       string
	lastWrapWidth int
	rendered      string
	width         int
	height        int
	focused       bool
}

func (p *Preview) SetFocus(f bool) {
	p.focused = f
}

func (p *Preview) SetContent(content string) {
	p.content = content
	p.rendered = "" // invalidate cache
	p.rerender()
}

func (p *Preview) SetSize(w, h int) {
	if w <= 2 || h <= 2 {
		return
	}
	p.width = w
	p.height = h
	p.viewport.SetWidth(w - 2)
	p.viewport.SetHeight(h - 2)
}

func (p *Preview) rerender() {
	if p.content == "" || p.width == 0 {
		return
	}

	wrapWidth := p.width - 4
	if wrapWidth < 20 {
		wrapWidth = 20
	}

	if wrapWidth == p.lastWrapWidth && p.rendered != "" {
		return
	}

	rd, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(wrapWidth),
	)
	if err != nil {
		return
	}

	rendered, err := rd.Render(p.content)
	if err != nil {
		return
	}

	p.rendered = rendered
	p.lastWrapWidth = wrapWidth

	prevOffset := p.viewport.YOffset()
	p.viewport.SetContent(rendered)
	p.viewport.SetYOffset(prevOffset)
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

	borderColor := lipgloss.Color("238")
	if p.focused {
		borderColor = lipgloss.Color("139")
	}

	return lipgloss.NewStyle().
		Width(p.width).
		Height(p.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Render(p.viewport.View())
}

func NewPreview() *Preview {
	vp := viewport.New(
		viewport.WithHeight(0),
		viewport.WithWidth(0),
	)
	return &Preview{
		viewport: vp,
	}
}
