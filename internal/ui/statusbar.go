package ui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type StatusBar struct {
	width int

	filePath string
	dirty    bool

	line      int
	column    int
	lineCount int
	wordCount int
	status    string

	spinner  spinner.Model
	spinning bool
	spinText string
}

func (s *StatusBar) SetSize(width int) {
	s.width = width
}

func (s *StatusBar) SetFile(path string) {
	s.filePath = path
}

func (s *StatusBar) SetDirty(d bool) {
	s.dirty = d
}

func (s *StatusBar) SetCursor(line, column int) {
	s.line = line
	s.column = column
}

func (s *StatusBar) SetStats(lines, words int) {
	s.lineCount = lines
	s.wordCount = words
}

func (s *StatusBar) SetStatus(msg string) {
	s.status = msg
}

// ClearStatus returns a Cmd that clears the status after a delay.
// Call this from Update whenever SetStatus is called.
func ClearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return ClearStatusMsg{}
	})
}

func (s *StatusBar) StartSpinner(text string) tea.Cmd {
	s.spinning = true
	s.spinText = text
	s.status = ""
	return s.spinner.Tick
}

func (s *StatusBar) StopSpinner() {
	s.spinning = false
	s.spinText = ""
}

func (s *StatusBar) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case ClearStatusMsg:
		s.status = ""
		return nil
	}

	if s.spinning {
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return cmd
	}

	return nil
}

func (s *StatusBar) View() string {
	dirtyIndicator := ""
	if s.dirty {
		dirtyIndicator = "*"
	}

	left := fmt.Sprintf(" %s%s", filepath.Base(s.filePath), dirtyIndicator)

	var right string
	switch {
	case s.spinning:
		right = fmt.Sprintf(" %s %s ", s.spinner.View(), s.spinText)
	case s.status != "":
		right = fmt.Sprintf(" %s ", s.status)
	default:
		right = fmt.Sprintf(
			"Ln %d, Col %d | %d lines | %d words ",
			s.line+1, s.column+1, s.lineCount, s.wordCount,
		)
	}

	gap := s.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}

	return lipgloss.NewStyle().
		Width(s.width).
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("252")).
		Render(left + strings.Repeat(" ", gap) + right)
}

func NewStatusBar() *StatusBar {
	s := spinner.New(
		spinner.WithSpinner(spinner.MiniDot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("141"))),
	)
	return &StatusBar{spinner: s}
}
