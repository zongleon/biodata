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

type choicePage struct {
	Title       string
	Description string
	Options     []int
	Labels      []string
	Choice      int
	Chosen      bool
}

func NewChoicePage(title string, desc string, choices []string, dests []int) *choicePage {
	options := make([]int, len(dests))
	copy(options, dests)

	return &choicePage{
		Title:       title,
		Description: desc,
		Options:     options,
		Labels:      choices,
		Choice:      0,
		Chosen:      false,
	}
}

// Choice page (menu)
func (page *choicePage) UpdatePage(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
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
			m.UpdateHistory(m.Page, page.Title)
			m.Page = page.Options[page.Choice]
		}
		// allow going back
		m.UpdateBack(msg)
	}
	return m, nil
}

func (page *choicePage) Page(m Model) string {
	choices := page.Labels

	tpl := page.Description + "\n"
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

func (page *choicePage) GetTitle() string {
	return page.Title
}
