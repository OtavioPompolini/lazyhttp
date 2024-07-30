package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type IWindow interface {
	Setup(w *Window)
	Update(w *Window)
	Size() (x, y, w, h int)
	Name() string
}

type Window struct {
	view   *gocui.View
	Window IWindow
	IsActive bool
}

func NewWindow(iw IWindow) *Window {
	return &Window{
		Window: iw,
		IsActive: false,
	}
}

func (v *Window) SetSelectedBgColor(col gocui.Attribute) {
	v.view.SelBgColor = col
}

func (v *Window) SetHightlight(b bool) {
	v.view.Highlight = b
}

func (v *Window) WriteLn(text string) {
	fmt.Fprintln(v.view, text)
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
