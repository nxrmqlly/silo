package ui

import (
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/nxrmqlly/silo/internal/fs"
)

type Sidebar struct {
	root    *fs.FileNode
	list    []*fs.FileNode // flattened render list
	cursor  int
	width   int
	height  int
	focused bool

	offset int

	mode      sidebarMode
	nameInput string // buffer since we only type one name at a time
}

type sidebarMode int

const (
	modeNormal        sidebarMode = iota
	modeNaming                    // 'n' pressed, typing new filename
	modeConfirmDelete             // 'd' pressed, waiting for y/n
)

func (s *Sidebar) refreshRenderList() {
	s.list = flattenTree(s.root, 0)
}

func flattenTree(node *fs.FileNode, level int) []*fs.FileNode {
	lines := []*fs.FileNode{}

	lines = append(lines, node)
	if node.IsDir && node.Expanded {
		for _, child := range node.Children {
			lines = append(lines, flattenTree(child, level+1)...)
		}
	}

	return lines
}

func depth(node *fs.FileNode) int {
	d := 0
	for node.Parent != nil {
		d++
		node = node.Parent
	}
	return d
}

func (s *Sidebar) adjustScroll() {
	if s.cursor < s.offset {
		s.offset = s.cursor
	}

	if s.cursor >= s.offset+s.height {
		s.offset = s.cursor - s.height + 1
	}
}

func (s *Sidebar) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func (s *Sidebar) SetFocus(f bool) {
	s.focused = f
}

// returns the directory path to create the new files in.
// If cursor is on a dir, use that else use parent dir.
func (s *Sidebar) dirAtCursor() string {
	if len(s.list) == 0 {
		return s.root.Path
	}

	current := s.list[s.cursor]
	if current.IsDir {
		return current.Path
	}
	if current.Parent != nil {
		return current.Parent.Path
	}

	return s.root.Path
}

func (s *Sidebar) handleNormalKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {

	case "down", "j":
		if s.cursor < len(s.list)-1 {
			s.cursor++
			s.adjustScroll()
		}

	case "up", "k":
		if s.cursor > 0 {
			s.cursor--
			s.adjustScroll()
		}

	case "enter":
		current := s.list[s.cursor]
		if !current.IsDir {
			return func() tea.Msg { return FileSelectedMsg{Path: current.Path} }
		}
		current.Expanded = !current.Expanded
		s.refreshRenderList()

	case "n":
		s.mode = modeNaming
		s.nameInput = ""

	case "d":
		if len(s.list) > 0 {
			s.mode = modeConfirmDelete
		}
	}

	return nil
}

func (s *Sidebar) handleNamingKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		s.mode = modeNormal
		s.nameInput = ""

	case "enter":
		input := strings.TrimSpace(s.nameInput)
		s.mode = modeNormal
		s.nameInput = ""
		if input == "" {
			return nil
		}

		dir := s.dirAtCursor()

		fullPath := filepath.Join(dir, input)
		isDir := strings.HasSuffix(input, "/")

		var err error

		if isDir {
			fullPath = strings.TrimSuffix(fullPath, "/")
			err = fs.CreateDir(fullPath)
		} else {
			err = fs.CreateFile(fullPath)
		}

		if err != nil {
			return nil
		}

		return func() tea.Msg {
			return FileCreatedMsg{Path: fullPath, IsDir: isDir}
		}
	case "backspace":
		if len(s.nameInput) > 0 {
			s.nameInput = s.nameInput[:len(s.nameInput)-1] // trim last char
		}

	default:
		ch := msg.String()

		// ensure only printable chars
		if len(ch) == 1 {
			s.nameInput += ch
		}
	}

	return nil
}

func (s *Sidebar) handleConfirmDeleteKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "y":
		s.mode = modeNormal
		if len(s.list) == 0 {
			return nil
		}
		target := s.list[s.cursor]
		path := target.Path

		if err := fs.DeletePath(path); err != nil {
			return nil
		}

		if s.cursor > 0 {
			s.cursor--
		}

		return func() tea.Msg {
			return FileDeletedMsg{Path: path}
		}
	case "n", "esc":
		s.mode = modeNormal
	}

	return nil
}

func (s *Sidebar) Update(msg tea.Msg) tea.Cmd {
	if !s.focused {
		return nil
	}

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch s.mode {

		case modeNaming:
			return s.handleNamingKey(msg)

		case modeConfirmDelete:
			return s.handleConfirmDeleteKey(msg)

		case modeNormal:
			return s.handleNormalKey(msg)
		}

	case RefreshSidebarMsg:
		fs.RefreshNode(s.root)
		s.refreshRenderList()

		//clamp the cursor
		if s.cursor >= len(s.list) {
			s.cursor = len(s.list) - 1 //last idx
		}
		return nil

	}

	return nil
}

func (s *Sidebar) View() string {
	var lines []string

	visible := s.list
	if s.offset < len(visible) {
		visible = visible[s.offset:]
	}

	maxLines := s.height - 2 // account for border
	if s.mode == modeNaming || s.mode == modeConfirmDelete {
		maxLines-- // reserve one line for the prompt
	}

	for i, node := range visible {
		if i >= maxLines {
			break
		}

		idx := i + s.offset
		prefix := strings.Repeat("  ", depth(node))
		name := node.Name

		if node.IsDir {
			if node.Expanded {
				name = "▾ " + name
			} else {
				name = "▸ " + name
			}
		}

		line := prefix + name

		if idx == s.cursor {
			if s.focused {
				line = lipgloss.NewStyle().
					Background(lipgloss.Color("238")).
					Foreground(lipgloss.Color("252")).
					Render(line)
			} else {
				line = lipgloss.NewStyle().
					Foreground(lipgloss.Color("244")).
					Render(line)
			}
		}

		lines = append(lines, line)
	}

	//prompt line at bottom
	switch s.mode {
	case modeNaming:
		prompt := lipgloss.NewStyle().
			Foreground(lipgloss.Color("139")).
			Render("new: " + s.nameInput + "|")
		lines = append(lines, prompt)

	case modeConfirmDelete:
		target := ""
		if len(s.list) > 0 {
			target = s.list[s.cursor].Name
		}

		prompt:= lipgloss.NewStyle().
		Foreground(lipgloss.Color("160")).
		Render("delete " + target + "? (y/n)")

		lines = append(lines, prompt)
	}
	

	style := lipgloss.NewStyle().
		Width(s.width).
		Height(s.height).
		Padding(0, 0)

	// border based on active or not
	if s.focused {
		style = style.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("139"))

	} else {
		style = style.Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238"))
	}

	return style.Render(strings.Join(lines, "\n"))
}

func NewSidebar(root *fs.FileNode) *Sidebar {
	s := &Sidebar{root: root}
	s.refreshRenderList()
	return s
}
