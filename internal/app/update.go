package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/nxrmqlly/silo/internal/ui"
)

func (m CustomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.editor.SetSize(m.width, m.height)
		m.statusbar.SetSize(msg.Width)

		return m, nil

	case ui.SaveFileMsg:
		// todo: false info - save logic later

		m.statusbar.SetFile(msg.Path)
		m.statusbar.SetDirty(false)
		
		return m, nil
	}

	cmd := m.editor.Update(msg)

	m.statusbar.SetFile(m.editor.FilePath())
	m.statusbar.SetDirty(m.editor.IsDirty())
	
	return m, cmd
}