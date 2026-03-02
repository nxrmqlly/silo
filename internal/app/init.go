package app

import (
	tea "charm.land/bubbletea/v2"
)

func (m CustomModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
