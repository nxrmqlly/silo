package ui

type OpenFileMsg struct {
	Path string
}

type SaveFileMsg struct {
	Path    string
	Content string
}

type FileSelectedMsg struct {
	Path string
}

type FileCreatedMsg struct {
	Path  string
	IsDir bool
}

type FileDeletedMsg struct {
	Path string
}

type RefreshSidebarMsg struct{}

type AutosaveMsg struct{}

type ClearStatusMsg struct{}

type PreviewRenderedMsg struct {
	Content  string // original source, for cache validation
	Rendered string
}

