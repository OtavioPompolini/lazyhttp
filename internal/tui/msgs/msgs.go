package msgs

import (
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/OtavioPompolini/project-postman/internal/state"
)

// FocusTarget identifies which pane is focused.
type FocusTarget int

const (
	FocusCollections    FocusTarget = iota
	FocusRequests       FocusTarget = iota
	FocusRequestDetails FocusTarget = iota
	FocusResponse       FocusTarget = iota
	FocusVariables      FocusTarget = iota
	FocusDebugger       FocusTarget = iota
	FocusCreateRequest  FocusTarget = iota
)

// ModalID identifies which modal overlay is active.
type ModalID int

const (
	ModalNone          ModalID = iota
	ModalAlert         ModalID = iota
	ModalCreateRequest ModalID = iota
)

// EventBus bridge messages
type CollectionChangedMsg struct{ Event state.CollectionsChangedEvent }
type RequestChangedMsg struct{ Event state.RequestsChangedEvent }
type CollectionSelectedMsg struct{ Event state.CollectionSelectedEvent }
type AlertMsg struct{ Message string }

// Focus / navigation
type FocusMsg struct{ Target FocusTarget }
type FocusLostMsg struct{}

// Modal lifecycle
type OpenModalMsg struct{ Modal ModalID }
type CloseModalMsg struct{}

// Internal tick messages
type AlertTickMsg struct{ Deadline time.Time }
type DebuggerTickMsg struct{}

// FocusCmd returns a Cmd that sends a FocusMsg.
func FocusCmd(target FocusTarget) tea.Cmd {
	return func() tea.Msg {
		return FocusMsg{Target: target}
	}
}

// CloseModalCmd returns a Cmd that sends a CloseModalMsg.
func CloseModalCmd() tea.Cmd {
	return func() tea.Msg {
		return CloseModalMsg{}
	}
}
