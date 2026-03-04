package panes

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/tui/msgs"
)

type CollectionsPane struct {
	focused          bool
	width, height    int
	lastEvent        state.CollectionEvent
	collectionSystem *state.CollectionSystem
	eventBus         *state.EventBus
}

func NewCollectionsPane(st *state.State, eb *state.EventBus) CollectionsPane {
	return CollectionsPane{
		collectionSystem: st.CollectionSystem,
		eventBus:         eb,
	}
}

func (p CollectionsPane) Init() tea.Cmd { return nil }

func (p CollectionsPane) Update(msg tea.Msg) (CollectionsPane, tea.Cmd) {
	switch m := msg.(type) {
	case msgs.CollectionChangedMsg:
		p.lastEvent = m.Event

	case msgs.FocusLostMsg:
		p.focused = false

	case tea.KeyMsg:
		if !p.focused {
			break
		}
		switch m.String() {
		case "j":
			p.collectionSystem.SelectNext()
		case "k":
			p.collectionSystem.SelectPrev()
		case "J":
			p.collectionSystem.SwapPositionDown()
		case "K":
			p.collectionSystem.SwapPositionUp()
		case "enter":
			if len(p.lastEvent.Collections) > 0 {
				id := p.lastEvent.Collections[p.lastEvent.CurrPos].Id
				p.eventBus.Publish(state.Event{
					Type: state.InternalCollectionSelected,
					Data: id,
				})
			}
		case "[":
			return p, msgs.FocusCmd(msgs.FocusRequests)
		}
	}
	return p, nil
}

func (p CollectionsPane) View() string {
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
	for i, col := range p.lastEvent.Collections {
		line := col.Name
		if i == p.lastEvent.CurrPos && i == p.lastEvent.SelPos {
			line = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true).Render("> " + line)
		} else if i == p.lastEvent.CurrPos {
			line = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("> " + line)
		} else if i == p.lastEvent.SelPos {
			line = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("  " + line)
		} else {
			line = "  " + line
		}
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(line)
	}

	title := lipgloss.NewStyle().Bold(true).Render("Collections")
	content := fmt.Sprintf("%s\n%s", title, sb.String())
	return border.Render(content)
}

func (p *CollectionsPane) SetSize(w, h int) {
	p.width = w
	p.height = h
}

func (p *CollectionsPane) SetFocused(f bool) {
	p.focused = f
}
