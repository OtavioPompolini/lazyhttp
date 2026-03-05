package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/tui/msgs"
	"github.com/OtavioPompolini/project-postman/internal/tui/panes"
)

type RootModel struct {
	width, height int
	focus         FocusManager
	activeModal   msgs.ModalID

	alertMessage string

	collections    panes.CollectionsPane
	requests       panes.RequestsPane
	requestDetails panes.RequestDetailsPane
	response       panes.ResponsePane
	variables      panes.VariablesPane
	debugger       panes.DebuggerPane
	createRequest  panes.CreateRequestPane
}

func NewRootModel(st *state.State, eb *state.EventBus) RootModel {
	m := RootModel{
		focus:          newFocusManager(),
		collections:    panes.NewCollectionsPane(st, eb),
		requests:       panes.NewRequestsPane(st),
		requestDetails: panes.NewRequestDetailsPane(st),
		response:       panes.NewResponsePane(),
		variables:      panes.NewVariablesPane(),
		debugger:       panes.NewDebuggerPane(),
		createRequest:  panes.NewCreateRequestPane(st),
	}
	m.applyFocus()
	return m
}

func (m RootModel) Init() tea.Cmd {
	return m.debugger.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch tm := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = tm.Width
		m.height = tm.Height
		m.distributeSize()
		return m, nil

	case tea.KeyMsg:
		// Global quit
		if tm.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Modal gets priority
		if m.activeModal == msgs.ModalCreateRequest {
			var cmd tea.Cmd
			m.createRequest, cmd = m.createRequest.Update(msg)
			return m, cmd
		}

		// Route to focused pane
		return m.routeKey(msg)

	case msgs.CollectionChangedMsg:
		var cmd tea.Cmd
		m.collections, cmd = m.collections.Update(msg)
		return m, cmd

	case msgs.RequestChangedMsg:
		var (
			cmd1, cmd2, cmd3 tea.Cmd
		)
		m.requests, cmd1 = m.requests.Update(msg)
		m.requestDetails, cmd2 = m.requestDetails.Update(msg)
		m.response, cmd3 = m.response.Update(msg)
		return m, tea.Batch(cmd1, cmd2, cmd3)

	case msgs.AlertMsg:
		m.alertMessage = tm.Message
		m.activeModal = msgs.ModalAlert
		deadline := time.Now().Add(2 * time.Second)
		return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
			return msgs.AlertTickMsg{Deadline: deadline}
		})

	case msgs.AlertTickMsg:
		if m.activeModal == msgs.ModalAlert {
			m.activeModal = msgs.ModalNone
			m.alertMessage = ""
		}
		return m, nil

	case msgs.DebuggerTickMsg:
		var cmd tea.Cmd
		m.debugger, cmd = m.debugger.Update(msg)
		return m, cmd

	case msgs.FocusMsg:
		m.loseFocus()
		m.focus.MoveTo(tm.Target)
		m.applyFocus()
		m.distributeSize()
		return m, nil

	case msgs.OpenModalMsg:
		m.activeModal = tm.Modal
		var cmd tea.Cmd
		m.createRequest, cmd = m.createRequest.Update(msg)
		return m, cmd

	case msgs.CloseModalMsg:
		m.activeModal = msgs.ModalNone
		return m, nil
	}

	return m, nil
}

func (m *RootModel) routeKey(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focus.Current() {
	case msgs.FocusCollections:
		m.collections, cmd = m.collections.Update(msg)
	case msgs.FocusRequests:
		m.requests, cmd = m.requests.Update(msg)
	case msgs.FocusRequestDetails:
		m.requestDetails, cmd = m.requestDetails.Update(msg)
	case msgs.FocusResponse:
		m.response, cmd = m.response.Update(msg)
	case msgs.FocusVariables:
		m.variables, cmd = m.variables.Update(msg)
	case msgs.FocusDebugger:
		m.debugger, cmd = m.debugger.Update(msg)
	}
	return m, cmd
}

// loseFocus sends FocusLostMsg to the currently focused pane.
func (m *RootModel) loseFocus() {
	lostMsg := msgs.FocusLostMsg{}
	switch m.focus.Current() {
	case msgs.FocusCollections:
		m.collections, _ = m.collections.Update(lostMsg)
	case msgs.FocusRequests:
		m.requests, _ = m.requests.Update(lostMsg)
	case msgs.FocusRequestDetails:
		m.requestDetails, _ = m.requestDetails.Update(lostMsg)
	case msgs.FocusResponse:
		m.response, _ = m.response.Update(lostMsg)
	case msgs.FocusVariables:
		m.variables, _ = m.variables.Update(lostMsg)
	case msgs.FocusDebugger:
		m.debugger, _ = m.debugger.Update(lostMsg)
	}
}

// applyFocus marks the currently focused pane as focused.
func (m *RootModel) applyFocus() {
	m.collections.SetFocused(m.focus.Current() == msgs.FocusCollections)
	m.requests.SetFocused(m.focus.Current() == msgs.FocusRequests)
	m.requestDetails.SetFocused(m.focus.Current() == msgs.FocusRequestDetails)
	m.response.SetFocused(m.focus.Current() == msgs.FocusResponse)
	m.variables.SetFocused(m.focus.Current() == msgs.FocusVariables)
	m.debugger.SetFocused(m.focus.Current() == msgs.FocusDebugger)
}

func (m RootModel) leftColHeights() (collH, reqH, varH int) {
	debugH := m.height * 20 / 100
	mainH := m.height - debugH

	switch m.focus.Current() {
	case msgs.FocusCollections:
		collH = mainH * 55 / 100
		reqH = mainH * 25 / 100
	case msgs.FocusRequests:
		collH = mainH * 25 / 100
		reqH = mainH * 55 / 100
	default:
		collH = mainH * 45 / 100
		reqH = mainH * 30 / 100
	}
	varH = mainH - collH - reqH
	return
}

func (m *RootModel) distributeSize() {
	leftW := m.width * 20 / 100
	midW := m.width * 40 / 100
	rightW := m.width - leftW - midW

	debugH := m.height * 20 / 100
	mainH := m.height - debugH

	collH, reqH, varH := m.leftColHeights()

	m.collections.SetSize(leftW, collH)
	m.requests.SetSize(leftW, reqH)
	m.variables.SetSize(leftW, varH)
	m.requestDetails.SetSize(midW, mainH)
	m.response.SetSize(rightW, mainH)
	m.debugger.SetSize(m.width, debugH)
	m.createRequest.SetSize(m.width, m.height)
}

func (m RootModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	leftW := m.width * 20 / 100
	midW := m.width * 40 / 100
	rightW := m.width - leftW - midW
	debugH := m.height * 20 / 100
	mainH := m.height - debugH
	collH, reqH, varH := m.leftColHeights()

	m.collections.SetSize(leftW, collH)
	m.requests.SetSize(leftW, reqH)
	m.variables.SetSize(leftW, varH)
	m.requestDetails.SetSize(midW, mainH)
	m.response.SetSize(rightW, mainH)
	m.debugger.SetSize(m.width, debugH)

	leftCol := lipgloss.JoinVertical(lipgloss.Left,
		m.collections.View(),
		m.requests.View(),
		m.variables.View(),
	)

	mainRow := lipgloss.JoinHorizontal(lipgloss.Top,
		leftCol,
		m.requestDetails.View(),
		m.response.View(),
	)

	debugRow := m.debugger.View()
	base := lipgloss.JoinVertical(lipgloss.Left, mainRow, debugRow)

	// Overlay modals
	if m.activeModal == msgs.ModalAlert {
		overlay := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("1")).
			Padding(0, 2).
			Render("⚠  " + m.alertMessage)
		base = placeOverlay(m.width/2-20, m.height/2-2, overlay, base)
	}

	if m.activeModal == msgs.ModalCreateRequest {
		overlay := m.createRequest.View()
		base = placeOverlay(m.width/2-26, m.height/2-3, overlay, base)
	}

	return base
}

// placeOverlay places overlay text on top of background text at position (cx, cy).
// Simple implementation: split both by lines, overlay line by line.
func placeOverlay(cx, cy int, overlay, bg string) string {
	if cx < 0 {
		cx = 0
	}
	if cy < 0 {
		cy = 0
	}

	bgLines := splitLines(bg)
	ovLines := splitLines(overlay)

	for i, ovLine := range ovLines {
		row := cy + i
		if row >= len(bgLines) {
			break
		}
		bgRunes := []rune(bgLines[row])
		ovRunes := []rune(ovLine)

		// Pad bg line if needed
		for len(bgRunes) < cx {
			bgRunes = append(bgRunes, ' ')
		}

		// Replace characters starting at cx
		for j, r := range ovRunes {
			col := cx + j
			if col < len(bgRunes) {
				bgRunes[col] = r
			} else {
				bgRunes = append(bgRunes, r)
			}
		}
		bgLines[row] = string(bgRunes)
	}

	result := ""
	for i, line := range bgLines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

func splitLines(s string) []string {
	lines := []string{}
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	lines = append(lines, s[start:])
	return lines
}
