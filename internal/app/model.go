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
	preview   *ui.Preview
	welcome   *ui.Welcome

	isPreview bool
	isWelcome bool
}

const (
	FocusSidebar FocusMode = iota
	FocusRight
)

func NewSiloModel(notesDir string) *SiloModel {
	return &SiloModel{
		focus:     FocusRight,
		isPreview: false,
		isWelcome: true,

		editor:    ui.NewEditor(),
		statusbar: ui.NewStatusBar(),
		sidebar:   ui.NewSidebar(fs.BuildFileTree(notesDir)),
		preview:   ui.NewPreview(),
		welcome:   ui.NewWelcome(),
	}
}
