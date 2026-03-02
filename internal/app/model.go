package app

import (
	"github.com/nxrmqlly/silo/internal/ui"
	"github.com/nxrmqlly/silo/internal/fs"
)

type FocusMode int
type CustomModel struct {
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

func InitialModel() *CustomModel {

	return &CustomModel{
		focus:     FocusEditor,
		editor:    ui.NewEditor(),
		statusbar: ui.NewStatusBar(),
		sidebar:   ui.NewSidebar(fs.BuildFileTree("./.silo-test/notes")),
	}
}
