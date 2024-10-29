package internal

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// styles
var (
	mainStyle = lipgloss.NewStyle().MarginLeft(2)
)

type Page interface {
	UpdatePage(tea.Msg, Model) (tea.Model, tea.Cmd)
	Page(Model) string
}

type Model struct {
	Pages         map[int]Page
	PreviousPages []int
	Page          int
	Quitting      bool
	Keys          keyMap
	ShowHelp      bool
	Help          help.Model
	Height        int
	Width         int
}

// keybindings
// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Back  key.Binding
	Enter key.Binding
	Help  key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Back, k.Enter, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down}, // these are columns
		{k.Left, k.Right},
		{k.Back, k.Enter},
		{k.Help, k.Quit},
	}
}

var Keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("bksp", "previous page"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// model.Model functions
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Help.Width = msg.Width
		m.Height = msg.Height
		m.Width = msg.Width

	case tea.KeyMsg:
		// always quit no matter what screen
		if key.Matches(msg, m.Keys.Quit) {
			m.Quitting = true
			return m, tea.Quit
		}
		// always show the extended help no matter what screen
		if key.Matches(msg, m.Keys.Help) {
			m.Help.ShowAll = !m.Help.ShowAll
		}

	}
	// update the page
	return m.Pages[m.Page].UpdatePage(msg, m)
}

func (m Model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	s = m.Pages[m.Page].Page(m)

	// display help at the bottom-ish
	helpView := m.Help.View(m.Keys)
	height := m.Height - strings.Count(s, "\n") - 22
	height = max(height, 0)
	helpView = strings.Repeat("\n", height) + helpView

	if m.ShowHelp {
		return mainStyle.Render("\n" + s + helpView)
	}
	return mainStyle.Render("\n" + s)
}

func (m *Model) UpdateBack(msg tea.KeyMsg) {
	// helper function, allows going back
	if key.Matches(msg, m.Keys.Back) {
		m.Page = m.PreviousPages[len(m.PreviousPages)-1]
		m.PreviousPages = m.PreviousPages[:len(m.PreviousPages)-1]
	}
}
