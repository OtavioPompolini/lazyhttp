package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type IView interface {
	Setup(g *gocui.Gui, v *gocui.View)
	Update(g *gocui.Gui, v *gocui.View)
	Size() (x, y, w, h int)
	Name() string
	SetKeybindings(ui UI) error
}

type View struct {
	Window IView
	view   *gocui.View
}

func NewView(v IView) (*View, error) {
	view := &View{
		Window: v,
	}

	return view, nil
}

func (v *View) Layout(g *gocui.Gui) error {
	x, y, w, h := v.Window.Size()
	newView, err := g.SetView(v.Window.Name(), x, y, w, h)
	if err != nil {
		if err == gocui.ErrUnknownView {
			v.setView(newView)
			v.Window.Setup(g, newView)
			return nil
		}
		return err
	}

	v.Window.Update(g, newView)

	return nil
}

func (v *View) SetSelectedBgColor(col gocui.Attribute) {
	v.view.SelBgColor = col
}

func (v *View) SetHightlight(b bool) {
	v.view.Highlight = b
}

func (v *View) WriteLn(text string) {
	fmt.Fprintln(v.view, text)
}

func (v *View) setView(newView *gocui.View) {
	v.view = newView
}
