package panes

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/tui/msgs"
	"github.com/OtavioPompolini/project-postman/internal/types"
)

type vimMode int

const (
	vimNormal vimMode = iota
	vimInsert vimMode = iota
)

type RequestDetailsPane struct {
	focused        bool
	width, height  int
	mode           vimMode
	textarea       textarea.Model
	currentRequest *types.Request
	requestSystem  *state.RequestManager
}

func NewRequestDetailsPane(st *state.State) RequestDetailsPane {
	ta := textarea.New()
	ta.Placeholder = "Request body..."
	ta.ShowLineNumbers = false
	ta.SetWidth(40)
	ta.SetHeight(20)

	return RequestDetailsPane{
		textarea:      ta,
		requestSystem: st.RequestManager,
		mode:          vimNormal,
	}
}

func (p RequestDetailsPane) Init() tea.Cmd { return nil }

func (p RequestDetailsPane) Update(msg tea.Msg) (RequestDetailsPane, tea.Cmd) {
	switch m := msg.(type) {
	case msgs.RequestChangedMsg:
		if !p.focused {
			p.currentRequest = nil
			if len(m.Event.Requests) > 0 && m.Event.Cursor >= 0 && m.Event.Cursor < len(m.Event.Requests) {
				p.currentRequest = m.Event.Requests[m.Event.Cursor]
			}
			body := ""
			if p.currentRequest != nil {
				body = p.currentRequest.Body
			}
			p.textarea.SetValue(body)
		}

	case msgs.FocusLostMsg:
		p.focused = false
		p.mode = vimNormal
		p.textarea.Blur()
		if p.currentRequest != nil {
			p.requestSystem.Update(&types.Request{
				Id:   p.currentRequest.Id,
				Body: p.textarea.Value(),
			})
		}

	case tea.KeyMsg:
		if !p.focused {
			break
		}
		switch p.mode {
		case vimNormal:
			switch m.String() {
			case "i":
				p.mode = vimInsert
				p.textarea.Focus()
				return p, textarea.Blink
			case "esc":
				return p, msgs.FocusCmd(msgs.FocusRequests)
			}
		case vimInsert:
			switch m.String() {
			case "esc":
				p.mode = vimNormal
				p.textarea.Blur()
			default:
				var cmd tea.Cmd
				p.textarea, cmd = p.textarea.Update(msg)
				return p, cmd
			}
		}
	}
	return p, nil
}

func (p RequestDetailsPane) View() string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(p.width - 2).
		Height(p.height - 2)

	if p.focused {
		border = border.BorderForeground(lipgloss.Color("2"))
	} else {
		border = border.BorderForeground(lipgloss.Color("8"))
	}

	modeStr := "NORMAL"
	if p.mode == vimInsert {
		modeStr = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("INSERT")
	}
	title := lipgloss.NewStyle().Bold(true).Render("Request Details") + "  [" + modeStr + "]"

	p.textarea.SetWidth(p.width - 4)
	p.textarea.SetHeight(p.height - 5)

	content := title + "\n" + p.textarea.View()
	return border.Render(content)
}

func (p *RequestDetailsPane) SetSize(w, h int) {
	p.width = w
	p.height = h
}

func (p *RequestDetailsPane) SetFocused(f bool) {
	p.focused = f
	if !f {
		p.mode = vimNormal
		p.textarea.Blur()
	}
}
