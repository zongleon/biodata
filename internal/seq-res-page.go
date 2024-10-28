package internal

import tea "github.com/charmbracelet/bubbletea"

type seqResPage struct {
	Data GBSeq
}

func newSeqResPage(data GBSeq) *seqResPage {
	return &seqResPage{
		Data: data,
	}
}

// UpdatePage implements page.
func (page *seqResPage) UpdatePage(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	m.ShowHelp = true
	return m, nil
}

// Page implements page.
func (page *seqResPage) Page(m Model) string {
	return "\n" + page.Data.PrettyPrint()
}
