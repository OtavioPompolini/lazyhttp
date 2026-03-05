package panes

import (
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

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
	pendingG       bool
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

	case tea.KeyPressMsg:
		if !p.focused {
			break
		}
		switch p.mode {
		case vimNormal:
			key := m.String()
			if key != "g" {
				p.pendingG = false
			}
			switch key {
			case "i", "a":
				p.mode = vimInsert
				return p, textarea.Blink
			case "I":
				p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl})
				p.mode = vimInsert
				return p, textarea.Blink
			case "A":
				p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: 'e', Mod: tea.ModCtrl})
				p.mode = vimInsert
				return p, textarea.Blink
			case "o":
				p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: 'e', Mod: tea.ModCtrl})
				p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
				p.mode = vimInsert
				return p, textarea.Blink
			case "O":
				p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl})
				p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
				p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: tea.KeyUp})
				p.mode = vimInsert
				return p, textarea.Blink
			case "esc":
				return p, msgs.FocusCmd(msgs.FocusRequests)
			default:
				if synth, ok := motionToKey(key, &p.pendingG); ok {
					var cmd tea.Cmd
					p.textarea, cmd = p.textarea.Update(synth)
					return p, cmd
				}
			}
		case vimInsert:
			switch m.String() {
			case "esc":
				p.mode = vimNormal
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
	if f {
		p.mode = vimNormal
		p.pendingG = false
		p.textarea.Focus()
	} else {
		p.mode = vimNormal
		p.textarea.Blur()
	}
}

// motionToKey translates a vim normal-mode motion key into a synthetic
// tea.KeyPressMsg the textarea can process. Returns (msg, true) if a motion
// was produced, or (_, false) if the key was consumed but produced no motion
// (e.g. the first 'g' of 'gg').
func motionToKey(key string, pendingG *bool) (tea.KeyPressMsg, bool) {
	switch key {
	case "h":
		return tea.KeyPressMsg{Code: tea.KeyLeft}, true
	case "l":
		return tea.KeyPressMsg{Code: tea.KeyRight}, true
	case "j":
		return tea.KeyPressMsg{Code: tea.KeyDown}, true
	case "k":
		return tea.KeyPressMsg{Code: tea.KeyUp}, true
	case "w", "e", "E":
		return tea.KeyPressMsg{Code: tea.KeyRight, Mod: tea.ModAlt}, true
	case "b", "B":
		return tea.KeyPressMsg{Code: tea.KeyLeft, Mod: tea.ModAlt}, true
	case "0":
		return tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl}, true
	case "$":
		return tea.KeyPressMsg{Code: 'e', Mod: tea.ModCtrl}, true
	case "G":
		return tea.KeyPressMsg{Code: tea.KeyEnd, Mod: tea.ModCtrl}, true
	case "g":
		if *pendingG {
			*pendingG = false
			return tea.KeyPressMsg{Code: tea.KeyHome, Mod: tea.ModCtrl}, true
		}
		*pendingG = true
		return tea.KeyPressMsg{}, false
	}
	return tea.KeyPressMsg{}, false
}
