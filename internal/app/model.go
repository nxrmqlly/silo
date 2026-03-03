package app

import (
	"github.com/nxrmqlly/silo/internal/fs"
	"github.com/nxrmqlly/silo/internal/ui"
)

type FocusMode int
type SiloModel struct {
	width  int
	height int

	focus FocusMode

	editor    *ui.Editor
	statusbar *ui.StatusBar
	sidebar   *ui.Sidebar
}

const (
	FocusSidebar FocusMode = iota
	FocusEditor
)

func NewSiloModel() *SiloModel {

	return &SiloModel{
		focus:     FocusEditor,
		editor:    ui.NewEditor(),
		statusbar: ui.NewStatusBar(),
		sidebar:   ui.NewSidebar(fs.BuildFileTree("./.silo-test/notes")),
	}
}
