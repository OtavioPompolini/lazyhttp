package panes

import (
	"bytes"
	"log"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/OtavioPompolini/project-postman/internal/tui/msgs"
)

// syncBuffer is a thread-safe bytes.Buffer used as log output.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (sb *syncBuffer) Write(p []byte) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.Write(p)
}

func (sb *syncBuffer) String() string {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.String()
}

type DebuggerPane struct {
	focused       bool
	width, height int
	viewport      viewport.Model
	logBuf        *syncBuffer
}

func NewDebuggerPane() DebuggerPane {
	buf := &syncBuffer{}
	log.SetOutput(buf)
	log.SetPrefix("INFO: ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	vp := viewport.New(80, 5)
	return DebuggerPane{
		viewport: vp,
		logBuf:   buf,
	}
}

func DebuggerTick() tea.Cmd {
	return tea.Every(100*time.Millisecond, func(t time.Time) tea.Msg {
		return msgs.DebuggerTickMsg{}
	})
}

func (p DebuggerPane) Init() tea.Cmd { return DebuggerTick() }

func (p DebuggerPane) Update(msg tea.Msg) (DebuggerPane, tea.Cmd) {
	switch msg.(type) {
	case msgs.DebuggerTickMsg:
		content := p.logBuf.String()
		p.viewport.SetContent(content)
		p.viewport.GotoBottom()
		return p, DebuggerTick()

	case msgs.FocusLostMsg:
		p.focused = false

	case tea.KeyMsg:
		if p.focused {
			var cmd tea.Cmd
			p.viewport, cmd = p.viewport.Update(msg)
			return p, cmd
		}
	}
	return p, nil
}

func (p DebuggerPane) View() string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(p.width - 2).
		Height(p.height - 2)

	if p.focused {
		border = border.BorderForeground(lipgloss.Color("2"))
	} else {
		border = border.BorderForeground(lipgloss.Color("8"))
	}

	title := lipgloss.NewStyle().Bold(true).Render("Debugger")
	p.viewport.Width = p.width - 4
	p.viewport.Height = p.height - 4

	content := title + "\n" + p.viewport.View()
	return border.Render(content)
}

func (p *DebuggerPane) SetSize(w, h int) {
	p.width = w
	p.height = h
	p.viewport.Width = w - 4
	p.viewport.Height = h - 4
}

func (p *DebuggerPane) SetFocused(f bool) {
	p.focused = f
}
