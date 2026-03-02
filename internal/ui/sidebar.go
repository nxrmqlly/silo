package ui

import (
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
}

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

func (s *Sidebar) Update(msg tea.Msg) tea.Cmd {
	if !s.focused {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
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
				return func() tea.Msg {
					return FileSelectedMsg{Path: current.Path}
				}
			}
			// toggle expand
			current.Expanded = !current.Expanded
			s.refreshRenderList()

		}
	}

	return nil
}

func (s *Sidebar) View() string {
	var lines []string

	for idx, node := range s.list {
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

	style := lipgloss.NewStyle().
		Width(s.width).
		Height(s.height).
		Padding(0, 0)

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
