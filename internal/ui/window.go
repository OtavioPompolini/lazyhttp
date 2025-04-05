package ui

import (
	"fmt"
	"io"

	"github.com/awesome-gocui/gocui"
)

const (
	FIXED    = "FIXED"
	RELATIVE = "RELATIVE"
)

type position string

type windowCoord struct {
	position position
	coord    int
}

func newWindowCoord(c int, pos position) windowCoord {
	return windowCoord{
		position: pos,
		coord:    c,
	}
}

type WindowPosition struct {
	x windowCoord
	y windowCoord
	h windowCoord
	w windowCoord
}

func NewWindowPosition(
	x, y, w, h int,
	posXa, posYa, posXb, posYb string,
) WindowPosition {
	return WindowPosition{
		x: newWindowCoord(x, position(posXa)),
		y: newWindowCoord(y, position(posYa)),
		w: newWindowCoord(w, position(posXb)),
		h: newWindowCoord(h, position(posYb)),
	}
}

// Size and name should be a Window attribute not a IWindow
type IWindow interface {
	Setup(ui *UI, w *Window)
	Update(ui UI, w Window)
	OnSelect(ui UI, w Window) error
	OnDeselect(ui UI, w Window) error
	Size() WindowPosition
	Name() string
	SetKeybindings(ui *UI, w *Window) error
	ReloadContent(ui *UI, w *Window)
}

type Window struct {
	view   *gocui.View
	Window IWindow
	// name           string
	isActive bool
}

// TODO: Builder pattern
func NewWindow(iw IWindow, ia bool) *Window {
	return &Window{
		// name:           name,
		Window:   iw,
		isActive: ia,
	}
}

// func (w *Window) Name() string {
// 	return w.name
// }

func (w *Window) EnableKeybindingOnEdit(b bool) {
	w.view.KeybindOnEdit = b
}

func (w *Window) IsActive() bool {
	return w.isActive
}

func (w *Window) OpenWindow() {
	w.isActive = true
}

func (w *Window) AutoScroll() {
	w.view.Autoscroll = true
}

func (v *Window) SetVimEditor() {
	v.view.Editor = NewVimEditor()
}

func (v *Window) SetEditor(e gocui.Editor) {
	v.view.Editor = e
}

func (v *Window) SetSelectedBgColor(col gocui.Attribute) {
	v.view.SelBgColor = col
}

func (v *Window) SetTitle(title string) {
	v.view.Title = title
}

func (v *Window) Wrap(b bool) {
	v.view.Wrap = b
}

func (v *Window) SetHightlight(b bool) {
	v.view.Highlight = b
}

// func (v *Window) SetFgColor(b bool) {
// 	v.view.FgColor = gocui.ColorRed
// }

func (v *Window) WriteLn(text string) {
	fmt.Fprint(v.view, "\n"+text)
}

func (v *Window) Write(text string) {
	fmt.Fprint(v.view, text)
}

func (v *Window) WriteHighlight(text string) {
	v.Write(text)
}
func (v *Window) WriteFunc(f func(wr io.Writer) error) error {
	return f(v.view)
}

func (v *Window) WriteLines(text []string) {
	for i, t := range text {
		if i < len(text)-1 {
			fmt.Fprintln(v.view, t)
		} else {
			fmt.Fprint(v.view, t)
		}
	}
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

func (v *Window) MoveCursorHalfWindowDown() {
	// _, _, _, d := v.Window.Size()
	d := 10
	v.view.MoveCursor(0, d/2)
}

func (v *Window) MoveCursorHalfWindowUp() {
	d := 10
	// _, _, _, d := v.Window.Size()
	v.view.MoveCursor(0, -d/2)
}

func (v *Window) SetCursor(x, y int) error {
	return v.view.SetCursor(x, y)
}
