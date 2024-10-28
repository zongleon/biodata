package internal

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Background(lipgloss.Color(""))
	selectedCardStyle = cardStyle.BorderBackground(lipgloss.Color("211"))
)

type ChoicePage struct {
	Title   string
	Options []int
	Labels  []string
	Choice  int
	Chosen  bool
}

func NewChoicePage(title string, choices []string, dests []int) *ChoicePage {
	options := make([]int, len(dests))
	copy(options, dests)

	return &ChoicePage{
		Title:   title,
		Options: options,
		Labels:  choices,
		Choice:  0,
		Chosen:  false,
	}
}

// Choice page (menu)
func (page *ChoicePage) UpdatePage(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Up):
			page.Choice--
			if page.Choice < 0 {
				page.Choice = 0
			}
		case key.Matches(msg, m.Keys.Down):
			page.Choice++
			if page.Choice > len(page.Labels)-1 {
				page.Choice = len(page.Labels) - 1
			}
		case key.Matches(msg, m.Keys.Left):
			// TODO
		case key.Matches(msg, m.Keys.Right):
			// TODO
		case key.Matches(msg, m.Keys.Enter):
			page.Chosen = true
			m.PreviousPages = append(m.PreviousPages, m.Page)
			m.Page = page.Options[page.Choice]
		}
	}
	return m, nil
}

func (page *ChoicePage) Page(m Model) string {
	choices := page.Labels

	tpl := page.Title + "\n\n"
	tpl += "%s\n\n"

	choiceText := ""
	var style lipgloss.Style
	idx := 0
	for _, choice := range choices {
		if idx == page.Choice {
			style = selectedCardStyle
		} else {
			style = cardStyle
		}
		choiceText += style.Render(choice) + "\n"
		idx++
	}

	return fmt.Sprintf(tpl, choiceText)
}
