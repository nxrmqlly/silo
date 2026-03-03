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
