package internal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type seqResPage struct {
	Data     GBSeq
	Viewport viewport.Model
	Title    string
	Width    int
}

func NewSeqResPage(data GBSeq, title string, width, height int) *seqResPage {
	headerHeight := lipgloss.Height(headerView(title, width))
	footerHeight := lipgloss.Height(footerView(width))
	verticalMarginHeight := headerHeight + footerHeight

	v := viewport.New(width, height-verticalMarginHeight)
	v.YPosition = headerHeight
	v.SetContent(data.PrettyPrint())
	return &seqResPage{
		Data:     data,
		Viewport: v,
		Title:    title,
		Width:    width,
	}
}

func headerView(title string, width int) string {
	titleBox := titleStyle.Render(title)
	line := strings.Repeat("─", max(0, width-lipgloss.Width(titleBox)))
	return lipgloss.JoinHorizontal(lipgloss.Center, titleBox, line)
}

func footerView(width int) string {
	line := strings.Repeat("─", width)
	return line
}

// UpdatePage implements page.
func (page *seqResPage) UpdatePage(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.UpdateBack(msg)
	}

	var cmd tea.Cmd
	page.Viewport, cmd = page.Viewport.Update(msg)

	return m, cmd
}

// Page implements page.
func (page *seqResPage) Page(m Model) string {
	return fmt.Sprintf("%s\n%s\n%s", headerView(page.Title, page.Width), page.Viewport.View(), footerView(page.Width))
}

func (page *seqResPage) GetTitle() string {
	return page.Title
}
