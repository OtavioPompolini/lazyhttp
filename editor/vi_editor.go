package Editor

import (
	"github.com/jroimartin/gocui"
)

const (
	NORMAL_MODE       string = "insert"
	INSERT_MODE       string = "normal"
	VISUAL_MODE       string = "visual"
	VISUAL_LINE_MODE  string = "visual_line"
	VISUAL_BLOCK_MODE string = "visual_block"
)

type VimEditor struct {
	Mode string
}

func (ve *VimEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if ve.Mode == INSERT_MODE {
		ve.InsertMode(v, key, ch, mod)
	} else if ve.Mode == NORMAL_MODE {
		ve.NormalMode(v, key, ch, mod)
	}
}

func (ve *VimEditor) InsertMode(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case key == gocui.KeyEsc:
		ve.Mode = NORMAL_MODE
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyEnter:
		v.EditNewLine()

	//// THIS SHOULD BE BANNED FOR VIMOTION USERS XD
	case key == gocui.KeyArrowDown:
		v.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
	////

	// TODO: handle other keybindings...
}

func (ve *VimEditor) NormalMode(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch == 'i':
		ve.Mode = INSERT_MODE
	case ch == 'j':
		v.MoveCursor(0, 1, false)
	case ch == 'k':
		v.MoveCursor(0, -1, false)
	case ch == 'h':
		v.MoveCursor(-1, 0, false)
	case ch == 'l':
		v.MoveCursor(1, 0, false)
	}
	// TODO: handle other keybindings...
}
