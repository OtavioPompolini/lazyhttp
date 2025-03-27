package editor

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type BubbleEditor struct {
	textarea *textarea.Model
}

func NewBubbleEditor() *BubbleEditor {
	model := textarea.New()
	model.Prompt = ""
	model.ShowLineNumbers = false
	model.KeyMap = textarea.VimKeyMap // Enable Vim bindings
	return &BubbleEditor{textarea: &model}
}

func (e *BubbleEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	// Convert gocui input to Bubble Tea message
	var msg tea.Msg
	switch key {
	case gocui.KeyArrowUp:
		msg = tea.KeyMsg{Type: tea.KeyUp}
	case gocui.KeyArrowDown:
		msg = tea.KeyMsg{Type: tea.KeyDown}
	case gocui.KeyArrowLeft:
		msg = tea.KeyMsg{Type: tea.KeyLeft}
	case gocui.KeyArrowRight:
		msg = tea.KeyMsg{Type: tea.KeyRight}
	case gocui.KeyEnter:
		msg = tea.KeyMsg{Type: tea.KeyEnter}
	case gocui.KeyBackspace, gocui.KeyBackspace2:
		msg = tea.KeyMsg{Type: tea.KeyBackspace}
	case gocui.KeyDelete:
		msg = tea.KeyMsg{Type: tea.KeyDelete}
	case gocui.KeyEsc:
		msg = tea.KeyMsg{Type: tea.KeyEscape}
	case gocui.KeyCtrlC:
		msg = tea.KeyMsg{Type: tea.KeyCtrlC}
	default:
		if ch != 0 {
			msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}}
		} else {
			return
		}
	}

	// Update textarea
	e.textarea.HandleKey(msg)

	// Update view content
	v.Clear()
	fmt.Fprint(v, e.textarea.Value())

	// Update cursor position
	x, y := e.textarea.Cursor()
	v.SetCursor(x, y)
}
