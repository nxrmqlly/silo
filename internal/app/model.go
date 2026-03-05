package app

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/nxrmqlly/silo/internal/fs"
	"github.com/nxrmqlly/silo/internal/ui"
)

type SiloModel struct {
	width  int
	height int

	focus     FocusMode
	rightPane RightPane

	editor    *ui.Editor
	statusbar *ui.StatusBar
	sidebar   *ui.Sidebar
	preview   *ui.Preview
	welcome   *ui.Welcome
}

type FocusMode int
type RightPane int

const (
	FocusSidebar FocusMode = iota
	FocusRight
)

const (
	PaneWelcome RightPane = iota
	PaneEditor
	PanePreview
)

func (m *SiloModel) setSbStatus(msg string) tea.Cmd {
	m.statusbar.SetStatus(msg)
	return ui.ClearStatusAfter(2 * time.Second)
}

func NewSiloModel(notesDir string) *SiloModel {
	return &SiloModel{
		focus:     FocusRight,
		rightPane: PaneWelcome,

		editor:    ui.NewEditor(),
		statusbar: ui.NewStatusBar(),
		sidebar:   ui.NewSidebar(fs.BuildFileTree(notesDir)),
		preview:   ui.NewPreview(),
		welcome:   ui.NewWelcome(),
	}
}
