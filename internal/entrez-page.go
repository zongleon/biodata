package internal

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

type entrezPage struct {
	Title       string
	Description string
	Filter      string
	Input       textinput.Model
	Response    []GBSeq
	Results     list.Model
	Spinner     spinner.Model
	Loading     bool
	Received    bool
}

func NewEntrezPage(title, filter, desc string) *entrezPage {
	ti := textinput.New()
	ti.Placeholder = "Text query"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &entrezPage{
		Title:       title,
		Description: desc,
		Filter:      filter,
		Input:       ti,
		Spinner:     s,
		Loading:     false,
		Received:    false,
	}
}

func seqSummaryItem(id string, seq GBSeq) ListItem {
	var title string
	if len(seq.References) > 0 {
		title = seq.References[0].Title
	} else {
		title = seq.Organism
	}
	return ListItem{
		title: seq.Definition,
		desc:  id + " - " + seq.MolType + " - " + title,
	}
}

func SeqsToItems(ids []string, seqs []GBSeq) []list.Item {
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
	result []GBSeq
	items  []list.Item
	query  string
}

type errMsg struct {
	err error
}

func fetch(filter, query string) func() tea.Msg {
	return func() tea.Msg {
		// hit the query endpoint
		ids, q, err := SearchDBForQuery("nuccore", filter, query)
		if err != nil {
			return errMsg{err: err}
		}
		res, err := EFetch("nuccore", ids)
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
func (page *entrezPage) UpdatePage(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	// new pages start
	pageStart := 1000

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Enter):
			if page.Received {
				break
			}
			if !page.Loading {
				page.Loading = true
				return m, tea.Batch(
					page.Spinner.Tick,
					fetch(page.Filter, page.Input.Value()),
				)
			}
		case key.Matches(msg, m.Keys.Back):
			if page.Input.Value() == "" && !page.Received {
				m.UpdateBack(msg)
			}
			if page.Received {
				page.Received = false
				page.Loading = false
				page.Input.Reset()
			}
		}

	case entrezMsg:
		page.Loading = false
		page.Response = msg.result

		// generate pages for all responses
		for idx, res := range page.Response {
			m.Pages[pageStart+idx] = NewSeqResPage(res, res.PrimaryAccession, m.Width-20, m.Height-8)
		}

		// customize delegate
		d := list.NewDefaultDelegate()
		d.UpdateFunc = UpdateDelegate
		page.Results = list.New(msg.items, d, m.Width-20, m.Height-8)
		page.Results.Title = msg.query
		m.ShowHelp = false
		page.Received = true
	case errMsg:
		log.Fatalf("error searching or fetching in Entrez")

	case listSelectMsg:
		m.ShowHelp = false
		m.UpdateHistory(m.Page, page.Title)
		m.Page = pageStart + msg.Index

	case tea.WindowSizeMsg:
		if page.Received {
			page.Results.SetSize(msg.Width-20, msg.Height-8)
		}
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

type listSelectMsg struct {
	Index int
}

func UpdateDelegate(msg tea.Msg, m *list.Model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.SelectedItem() == nil {
				return nil
			}
			selectedIndex := m.Index()
			return func() tea.Msg {
				return listSelectMsg{
					Index: selectedIndex,
				}
			}
		}
	}
	return nil
}

func (page *entrezPage) Page(m Model) string {
	p := page.Description + "\n"
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

func (page *entrezPage) GetTitle() string {
	return page.Title
}
