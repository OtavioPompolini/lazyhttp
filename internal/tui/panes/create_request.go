package panes

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/tui/msgs"
)

type CreateRequestPane struct {
	width, height int
	input         textinput.Model
	requestSystem *state.RequestManager
}

func NewCreateRequestPane(st *state.State) CreateRequestPane {
	ti := textinput.New()
	ti.Placeholder = "Request name..."
	ti.Focus()
	ti.CharLimit = 128
	return CreateRequestPane{
		input:         ti,
		requestSystem: st.RequestManager,
	}
}

func (p CreateRequestPane) Init() tea.Cmd { return textinput.Blink }

func (p CreateRequestPane) Update(msg tea.Msg) (CreateRequestPane, tea.Cmd) {
	switch m := msg.(type) {
	case msgs.OpenModalMsg:
		if m.Modal == msgs.ModalCreateRequest {
			p.input.SetValue("")
			p.input.Focus()
		}

	case tea.KeyPressMsg:
		switch m.String() {
		case "enter":
			name := p.input.Value()
			p.input.SetValue("")
			if name != "" {
				rs := p.requestSystem
				return p, tea.Batch(
					func() tea.Msg { rs.Create(name); return nil },
					msgs.CloseModalCmd(),
				)
			}
			return p, msgs.CloseModalCmd()
		case "esc":
			p.input.SetValue("")
			return p, msgs.CloseModalCmd()
		default:
			var cmd tea.Cmd
			p.input, cmd = p.input.Update(msg)
			return p, cmd
		}
	}
	return p, nil
}

func (p CreateRequestPane) View() string {
	modalW := 50
	modalH := 5

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("5")).
		Width(modalW - 2).
		Height(modalH - 2).
		Padding(0, 1)

	title := lipgloss.NewStyle().Bold(true).Render("New request name:")
	content := title + "\n" + p.input.View()
	return box.Render(content)
}

func (p *CreateRequestPane) SetSize(w, h int) {
	p.width = w
	p.height = h
}
