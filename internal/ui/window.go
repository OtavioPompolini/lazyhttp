package ui

import (
	"fmt"
	"io"
	"log"

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

type IWindow interface {
	Setup(ui *UI)
	Update(ui UI)
	SetKeybindings(ui *UI) error
	Size() WindowPosition
	Name() string

	//Still dont know if those methods are required in the interface
	OnSelect(ui UI, w Window) error
	OnDeselect(ui UI, w Window) error
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

func (v *Window) Cursor() (x, y int) {
	return v.view.Cursor()
}

func (v *Window) CursorDown() {
	v.view.MoveCursor(0, 1)
}

func (v *Window) CursorUp() {
	v.view.MoveCursor(0, -1)
}

func (v *Window) CursorLeft() {
	v.view.MoveCursor(-1, 0)
}

func (v *Window) CursorRight() {
	v.view.MoveCursor(1, 0)
}

func (v *Window) HighlightLine(y int, b bool) {
	v.view.SetHighlight(y, b)
}

func (v *Window) Line(y int) string {
	s, err := v.view.Line(y)
	if err != nil {
		log.Println("Error getting line content from view buffer")
		return ""
	}

	return s
}
