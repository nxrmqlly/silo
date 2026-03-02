package ui

type OpenFileMsg struct {
	Path string
}
type SaveFileMsg struct {
	Path    string
	Content string
}
