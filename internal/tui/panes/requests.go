package panes

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/tui/msgs"
)

type RequestsPane struct {
	focused       bool
	width, height int
	lastEvent     state.RequestsChangedEvent
	requestSystem *state.RequestManager
}

func NewRequestsPane(st *state.State) RequestsPane {
	return RequestsPane{
		requestSystem: st.RequestManager,
	}
}

func (p RequestsPane) Init() tea.Cmd { return nil }

func (p RequestsPane) Update(msg tea.Msg) (RequestsPane, tea.Cmd) {
	switch m := msg.(type) {
	case msgs.RequestChangedMsg:
		p.lastEvent = m.Event

	case msgs.FocusLostMsg:
		p.focused = false

	case tea.KeyMsg:
		if !p.focused {
			break
		}
		rs := p.requestSystem
		switch m.String() {
		case "j":
			return p, func() tea.Msg { rs.SelectNext(); return nil }
		case "k":
			return p, func() tea.Msg { rs.SelectPrev(); return nil }
		case "enter":
			return p, msgs.FocusCmd(msgs.FocusRequestDetails)
		case "1":
			return p, msgs.FocusCmd(msgs.FocusCollections)
		case "r":
			return p, msgs.FocusCmd(msgs.FocusResponse)
		case "n":
			return p, func() tea.Msg {
				return msgs.OpenModalMsg{Modal: msgs.ModalCreateRequest}
			}
		}
	}
	return p, nil
}

func (p RequestsPane) View() string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(p.width - 2).
		Height(p.height - 2)

	if p.focused {
		border = border.BorderForeground(lipgloss.Color("2"))
	} else {
		border = border.BorderForeground(lipgloss.Color("8"))
	}

	var sb strings.Builder
	for i, req := range p.lastEvent.Requests {
		line := req.Name
		if i == p.lastEvent.Cursor {
			line = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("> " + line)
		} else {
			line = "  " + line
		}
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(line)
	}

	title := lipgloss.NewStyle().Bold(true).Render("Requests")
	content := title + "\n" + sb.String()
	return border.Render(content)
}

func (p *RequestsPane) SetSize(w, h int) {
	p.width = w
	p.height = h
}

func (p *RequestsPane) SetFocused(f bool) {
	p.focused = f
}
