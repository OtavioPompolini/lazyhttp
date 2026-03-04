package panes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/OtavioPompolini/project-postman/internal/tui/msgs"
)

type VariablesPane struct {
	focused       bool
	width, height int
}

func NewVariablesPane() VariablesPane {
	return VariablesPane{}
}

func (p VariablesPane) Init() tea.Cmd { return nil }

func (p VariablesPane) Update(msg tea.Msg) (VariablesPane, tea.Cmd) {
	switch msg.(type) {
	case msgs.FocusLostMsg:
		p.focused = false
	}
	return p, nil
}

func (p VariablesPane) View() string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(p.width - 2).
		Height(p.height - 2)

	if p.focused {
		border = border.BorderForeground(lipgloss.Color("2"))
	} else {
		border = border.BorderForeground(lipgloss.Color("8"))
	}

	title := lipgloss.NewStyle().Bold(true).Render("Variables")
	return border.Render(title)
}

func (p *VariablesPane) SetSize(w, h int) {
	p.width = w
	p.height = h
}

func (p *VariablesPane) SetFocused(f bool) {
	p.focused = f
}
