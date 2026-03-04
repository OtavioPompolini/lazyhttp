package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/tui/msgs"
)

// Bridge wires EventBus callbacks into tea.Program.Send() calls.
// program.Send is goroutine-safe by design.
type Bridge struct {
	program *tea.Program
}

func NewBridge(p *tea.Program, eb *state.EventBus) *Bridge {
	b := &Bridge{program: p}

	eb.Subscribe(state.CollectionChanged, func(e state.Event) {
		b.program.Send(msgs.CollectionChangedMsg{Event: e.Data.(state.CollectionEvent)})
	})

	eb.Subscribe(state.CollectionSelected, func(e state.Event) {
		b.program.Send(msgs.CollectionSelectedMsg{Event: e.Data.(state.CollectionSelectedEvent)})
	})

	eb.Subscribe(state.RequestChanged, func(e state.Event) {
		b.program.Send(msgs.RequestChangedMsg{Event: e.Data.(state.RequestEvent)})
	})

	eb.Subscribe(state.AlertMessage, func(e state.Event) {
		b.program.Send(msgs.AlertMsg{Message: e.Data.(string)})
	})

	return b
}
