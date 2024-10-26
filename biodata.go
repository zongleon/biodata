package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/zongleon/biodata/entrez"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
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
	Back  key.Binding
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

// styles
var (
	mainStyle = lipgloss.NewStyle().MarginLeft(2)
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Background(lipgloss.Color(""))
	// selectedCardStyle = cardStyle.Background(lipgloss.Color("211")).
	selectedCardStyle = cardStyle.BorderBackground(lipgloss.Color("211"))
)

type page interface {
	UpdatePage(tea.Msg, model) (tea.Model, tea.Cmd)
	Page(model) string
}
type choicePage struct {
	Title   string
	Options []int
	Labels  []string
}
type entrezPage struct {
	Title    string
	Input    textinput.Model
	Response []entrez.GBSeq
	Results  list.Model
	Spinner  spinner.Model
	Loading  bool
	Received bool
}

func newChoicePage(title string, choices []string, dests []int) page {
	options := make([]int, len(dests))
	copy(options, dests)

	return &choicePage{
		Title:   title,
		Options: options,
		Labels:  choices,
	}
}

func newEntrezPage(title string) page {
	ti := textinput.New()
	ti.Placeholder = "Text query"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &entrezPage{
		Title:    title,
		Input:    ti,
		Spinner:  s,
		Loading:  false,
		Received: false,
	}
}

type model struct {
	Pages         map[int]page
	PreviousPages []int
	Page          int
	Choice        int
	Chosen        bool
	Quitting      bool
	Keys          keyMap
	Help          help.Model
	Height        int
	Width         int
}

func newModel() model {
	return model{
		Pages: map[int]page{
			// type, sub-type
			0: newChoicePage("What type of biological data?", []string{"DNA", "RNA", "Protein", "Literature"}, []int{1, 2, 3, 4}),
			1: newChoicePage("What sort of DNA data?", []string{"Genome", "Genes", "Variation"}, []int{5, 6, 7}),
			2: newChoicePage("What sort of RNA data?", []string{"Transcript", "Expression"}, []int{8}),
			3: newChoicePage("What sort of protein data?", []string{"Sequence", "Structure", "Interactions"}, []int{9, 10, 11}),
			4: newChoicePage("What sort of literature?", []string{"Keyword", "Author"}, []int{12, 13}),

			// databases
			5: newChoicePage("Here are some genomic databases. Choose one to see what type of queries you can make.",
				[]string{"GenBank", "RefSeq", "Ensembl", "UCSC Genome Browser"}, []int{14, 15, 16, 17}),

			// access
			14: newEntrezPage("GenBank is an archival NCBI dataset, containing all publicly submitted DNA sequences from individual labs and large-scale sequencing projects."),
		},
		PreviousPages: []int{},
		Page:          0,
		Choice:        0,
		Chosen:        false,
		Quitting:      false,
		Keys:          keys,
		Help:          help.New(),
		Height:        40,
		Width:         40,
	}
}

// Model functions
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		// update the page
	}
	return m.Pages[m.Page].UpdatePage(msg, m)
}

func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	s = m.Pages[m.Page].Page(m)

	// display help at the bottom-ish
	helpView := m.Help.View(m.Keys)
	height := m.Height - strings.Count(s, "\n") - 22
	height = max(height, 0)

	return mainStyle.Render("\n" + s + strings.Repeat("\n", height) + helpView)
}

// Choice page (menu)
func (page *choicePage) UpdatePage(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Up):
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case key.Matches(msg, m.Keys.Down):
			m.Choice++
			if m.Choice > len(page.Labels)-1 {
				m.Choice = len(page.Labels) - 1
			}
		case key.Matches(msg, m.Keys.Left):
			// TODO
		case key.Matches(msg, m.Keys.Right):
			// TODO
		case key.Matches(msg, m.Keys.Enter):
			m.Chosen = true
			m.PreviousPages = append(m.PreviousPages, m.Page)
			m.Page = page.Options[m.Choice]
		}
	}
	return m, nil
}

func (page *choicePage) Page(m model) string {
	choices := page.Labels

	tpl := page.Title + "\n\n"
	tpl += "%s\n\n"

	choiceText := ""
	var style lipgloss.Style
	idx := 0
	for _, choice := range choices {
		if idx == m.Choice {
			style = selectedCardStyle
		} else {
			style = cardStyle
		}
		choiceText += style.Render(choice) + "\n"
		idx++
	}

	return fmt.Sprintf(tpl, choiceText)
}

func seqSummaryItem(id string, seq entrez.GBSeq) ListItem {
	return ListItem{
		title: seq.Definition,
		desc:  id + " - " + seq.MolType + " - " + seq.References[0].Title,
	}
}

func SeqsToItems(ids []string, seqs []entrez.GBSeq) []list.Item {
	out := make([]list.Item, len(seqs))
	for idx, seq := range seqs {
		out[idx] = seqSummaryItem(ids[idx], seq)
	}
	return out
}

type ListItem struct {
	title, desc string
}

func (i ListItem) Title() string       { return i.title }
func (i ListItem) Description() string { return i.desc }
func (i ListItem) FilterValue() string { return i.title }

type entrezMsg struct {
	result []entrez.GBSeq
	items  []list.Item
	query  string
}

type errMsg struct {
	err error
}

func fetch(query string) func() tea.Msg {
	return func() tea.Msg {
		// hit the query endpoint
		ids, q, err := entrez.SearchDBForQuery("nuccore", query)
		if err != nil {
			return errMsg{err: err}
		}
		res, err := entrez.EFetch("nuccore", ids)
		if err != nil {
			return errMsg{err: err}
		}
		return entrezMsg{
			result: res,
			items:  SeqsToItems(ids, res),
			query:  q,
		}
	}
}

// Entrez page
func (page *entrezPage) UpdatePage(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Up):
			// TODO
		case key.Matches(msg, m.Keys.Down):
			// TODO
		case key.Matches(msg, m.Keys.Enter):
			if page.Received {
				break
			}
			if !page.Loading {
				page.Loading = true
				return m, tea.Batch(
					page.Spinner.Tick,
					fetch(page.Input.Value()),
				)
			}
		}
	case entrezMsg:
		page.Loading = false
		page.Response = msg.result
		page.Results = list.New(msg.items, list.NewDefaultDelegate(), m.Width-20, m.Height-8)
		page.Results.Title = msg.query
		m.Help = page.Results.Help
		page.Received = true
	case errMsg:

	case tea.WindowSizeMsg:
		page.Results.SetSize(msg.Width-20, msg.Height-8)
	}

	// update page
	var cmd tea.Cmd
	if page.Loading {
		page.Spinner, cmd = page.Spinner.Update(msg)
	} else if page.Received {
		page.Results, cmd = page.Results.Update(msg)
	} else {
		page.Input, cmd = page.Input.Update(msg)
	}

	return m, cmd
}

func (page *entrezPage) Page(m model) string {
	p := page.Title
	p += "\n\n"

	if page.Received {
		p += page.Results.View()
	} else if page.Loading {
		p += page.Spinner.View() + " Loading results of query ... "
	} else {
		p += page.Input.View()
	}

	p += "\n\n"
	return p
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}

}
