package panes

import (
	"bytes"
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/OtavioPompolini/project-postman/internal/tui/msgs"
	"github.com/OtavioPompolini/project-postman/internal/types"
	"github.com/OtavioPompolini/project-postman/internal/utils"
)

type ResponsePane struct {
	focused        bool
	width, height  int
	viewport       viewport.Model
	currentRequest *types.Request
}

func NewResponsePane() ResponsePane {
	vp := viewport.New(60, 20)
	return ResponsePane{viewport: vp}
}

func (p ResponsePane) Init() tea.Cmd { return nil }

func (p ResponsePane) Update(msg tea.Msg) (ResponsePane, tea.Cmd) {
	switch m := msg.(type) {
	case msgs.RequestChangedMsg:
		p.currentRequest = nil
		if len(m.Event.Requests) > 0 && m.Event.Cursor >= 0 && m.Event.Cursor < len(m.Event.Requests) {
			p.currentRequest = m.Event.Requests[m.Event.Cursor]
		}
		p.viewport.SetContent(p.buildContent())

	case msgs.FocusLostMsg:
		p.focused = false

	case tea.KeyMsg:
		if !p.focused {
			break
		}
		switch m.String() {
		case "f1":
			return p, msgs.FocusCmd(msgs.FocusRequests)
		default:
			var cmd tea.Cmd
			p.viewport, cmd = p.viewport.Update(msg)
			return p, cmd
		}
	}
	return p, nil
}

func (p ResponsePane) buildContent() string {
	if p.currentRequest == nil || len(p.currentRequest.ResponseHistory) == 0 {
		return ""
	}
	resp := p.currentRequest.ResponseHistory[0]
	var buf bytes.Buffer
	buf.WriteString(resp.Info)
	buf.WriteString("\n")
	if err := utils.FormatAndHighlight(resp.Body)(&buf); err != nil {
		log.Print("Error while highlighting response body:", err)
		buf.WriteString(resp.Body)
	}
	return buf.String()
}

func (p ResponsePane) View() string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(p.width - 2).
		Height(p.height - 2)

	if p.focused {
		border = border.BorderForeground(lipgloss.Color("2"))
	} else {
		border = border.BorderForeground(lipgloss.Color("8"))
	}

	title := lipgloss.NewStyle().Bold(true).Render("Response")
	p.viewport.Width = p.width - 4
	p.viewport.Height = p.height - 4

	content := title + "\n" + p.viewport.View()
	return border.Render(content)
}

func (p *ResponsePane) SetSize(w, h int) {
	p.width = w
	p.height = h
	p.viewport.Width = w - 4
	p.viewport.Height = h - 4
}

func (p *ResponsePane) SetFocused(f bool) {
	p.focused = f
}
