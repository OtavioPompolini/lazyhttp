package ui

import (
	"fmt"

	Editor "github.com/OtavioPompolini/project-postman/editor"
	"github.com/jroimartin/gocui"
)

type IWindow interface {
	Setup(w Window)
	Update(w Window)
	OnSelect() error
	OnDeselect() error
	Size() (x, y, w, h int)
	Name() string
	IsActive() bool
	SetKeybindings(ui *UI) error
}

type Window struct {
	view   *gocui.View
	Window IWindow
}

//TODO: Builder pattern
func NewWindow(iw IWindow) *Window {
	return &Window{
		Window: iw,
	}
}

func (v *Window) SetVimEditor() {
	v.view.Editor = &Editor.VimEditor{}
}

func (v *Window) SetSelectedBgColor(col gocui.Attribute) {
	v.view.SelBgColor = col
}

func (v *Window) SetTitle(title string) {
	v.view.Title = title
}

func (v *Window) SetHightlight(b bool) {
	v.view.Highlight = b
}

//FIX: WriteLn always put a '\n' at the end so it will have a new blank line
func (v *Window) WriteLn(text string) {
	fmt.Fprintln(v.view, text)
}

func (v *Window) WriteLines(text []string) {
	for i, t := range text {
		if (i < len(text)-1) {
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
