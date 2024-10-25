package main

// An example demonstrating an application with multiple views.
//
// Note that this example was produced before the Bubbles progress component
// was available (github.com/charmbracelet/bubbles/progress) and thus, we're
// implementing a progress bar from scratch here.

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// keybindings
// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Help  key.Binding
	Quit  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}

var keys = keyMap{
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

// General stuff for styling the view
var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	mainStyle    = lipgloss.NewStyle().MarginLeft(2)
)

type page struct {
	// UpdatePage func(msg tea.Msg, m model) (tea.Model, tea.Cmd)
	// Page       func(m model) string
	Options []string
}

type model struct {
	Pages    []page
	Choice   int
	Chosen   bool
	Quitting bool
	keys     keyMap
	help     help.Model
}

func newModel() model {
	return model{
		Pages: []page{
			{Options: []string{"DNA / RNA", "Protein", "Literature", "Proteins"}},
		},
		Choice:   0,
		Chosen:   false,
		Quitting: false,
		keys:     keys,
		help:     help.New(),
	}
}

// Model functions
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// always quit no matter what screen
		if key.Matches(msg, m.keys.Quit) {
			m.Quitting = true
			return m, tea.Quit
		}
		// update the choice screen
		if !m.Chosen {
			return updateChoices(msg, m)
		}
		// update the page
		// return m.Pages[m.Choice].UpdatePage(msg, m)

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	}
	return m, nil
}

func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	if !m.Chosen {
		s = choicesView(m)
	} else {
		// s = m.Pages[m.Choice].Page(m)
	}
	return mainStyle.Render("\n" + s + "\n\n")
}

// View functions

// View chooser
func updateChoices(msg tea.KeyMsg, m model) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		m.Choice--
		if m.Choice < 0 {
			m.Choice = 0
		}
	case key.Matches(msg, m.keys.Down):
		m.Choice++
		if m.Choice > len(m.Pages)-1 {
			m.Choice = len(m.Pages) - 1
		}

	case key.Matches(msg, m.keys.Left):
		// TODO
	case key.Matches(msg, m.keys.Right):
		// TODO
	case key.Matches(msg, m.keys.Enter):
		m.Chosen = true
	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll
	}
	return m, nil
}

// Choices

// select between various choices
func choicesView(m model) string {
	choices := m.Pages[m.Choice].Options

	tpl := "What biological data are you interested in?\n\n"
	tpl += "%s\n\n"

	// helpView := m.help.View(m.keys)
	// height := 8 - strings.Count(status, "\n") - strings.Count(helpView, "\n")

	tpl += m.help.View(m.keys)

	choiceText := ""
	for _, choice := range choices {
		choiceText += choice + "\n"
	}

	return fmt.Sprintf(tpl, choiceText)
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
