package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

// Size and name should be a Window attribute not a IWindow
type IWindow interface {
	Setup(ui UI, w Window)
	Update(ui UI, w Window)
	OnSelect(ui UI, w Window) error
	OnDeselect(ui UI, w Window) error
	Size() (x, y, w, h int)
	Name() string
	SetKeybindings(ui *UI) error
}

type Window struct {
	view     *gocui.View
	Window   IWindow
	isActive bool
}

// TODO: Builder pattern
func NewWindow(iw IWindow, ia bool) *Window {
	return &Window{
		Window:   iw,
		isActive: ia,
	}
}

func (w *Window) IsActive() bool {
	return w.isActive
}

func (w *Window) SwitchOnOff(b bool) {
	w.isActive = b
}

// func (v *Window) SetVimEditor() {
// 	v.view.Editor = &Editor.VimEditor{}
// }

func (v *Window) SetSelectedBgColor(col gocui.Attribute) {
	v.view.SelBgColor = col
}

func (v *Window) SetTitle(title string) {
	v.view.Title = title
}

func (v *Window) SetHightlight(b bool) {
	v.view.Highlight = b
}

func (v *Window) SetFgColor(b bool) {
	v.view.FgColor = gocui.ColorRed
}

// FIX: WriteLn always put a '\n' at the end so it will have a new blank line
func (v *Window) WriteLn(text string) {
	fmt.Fprintln(v.view, text)
}

func (v *Window) WriteLines(text []string) {
	for _, t := range text {
		fmt.Fprintln(v.view, t)
	}

	// for i, t := range text {
	// 	if i < len(text)-1 {
	// 		fmt.Fprintln(v.view, t)
	// 	} else {
	// 		fmt.Fprint(v.view, t)
	// 	}
	// }
}

func (v *Window) GetWindowContent() string {
	return v.view.Buffer()
}

func (v *Window) ClearWindow() {
	v.view.Clear()
}

func (v *Window) SetEditable(b bool) {
	v.view.Editable = b
}

func (v *Window) setView(newView *gocui.View) {
	v.view = newView
}

func (v *Window) IsTained() bool {
	return v.view.IsTainted()
}

func (v *Window) SetCursor(x, y int) error {
	return v.view.SetCursor(x, y)
}
