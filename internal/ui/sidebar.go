package ui

import (
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Sidebar struct {
	files   []string
	cursor  int
	offset  int // scroll offset
	width   int
	height  int
	focused bool
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

		case "j", "down":
			if s.cursor < len(s.files)-1 {
				s.cursor++
			}

		case "k", "up":
			if s.cursor > 0 {
				s.cursor--
			}

		case "enter":
			if len(s.files) == 0 {
				return nil
			}
			selected := s.files[s.cursor]
			return func() tea.Msg {
				return FileSelectedMsg{Path: selected}
			}
		}
	}

	s.adjustScroll()

	return nil
}

func (s *Sidebar) View() string {
	var out string

	// calc end before overflow
	end := s.offset + s.height
	if end > len(s.files) {
		end = len(s.files)
	}

	// loop offset to end
	for i := s.offset; i < end; i++ {
		file := filepath.Base(s.files[i])
		line := " " + file

		if i == s.cursor {
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

		out += line + "\n"
	}

	outStyle := lipgloss.NewStyle()
	return outStyle.Width(s.width).Render(out)
}

func NewSidebar(files []string) *Sidebar {
	return &Sidebar{
		files: files,
	}
}
