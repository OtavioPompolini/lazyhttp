package ui

import (
	"errors"
	"fmt"

	"github.com/jroimartin/gocui"
)

type RequestsWindow struct {
	name                string
	x, y                int
	OnOpenAddNewRequest func() error
}

func (w RequestsWindow) Name() string {
	return w.name
}

func NewRequestsWindow(ui *UI) *View {
	return &View{
		Window: &RequestsWindow{
			name: "RequestsWindow",
			x:    0,
			y:    0,
		},
	}
}

func (w *RequestsWindow) Setup(g *gocui.Gui, v *gocui.View) {
	g.SetCurrentView(w.name)
	v.SelBgColor = gocui.ColorRed
	v.Highlight = true
	// v.SetSelectedBgColor(gocui.ColorRed)
	// v.SetHightlight(true)
}

func (w *RequestsWindow) Update(g *gocui.Gui, v *gocui.View) {
	v.Clear()
	for i := 0; i < 10; i++ {
		fmt.Fprintln(v, "PUDIM")
	}
}

func (w *RequestsWindow) Size() (x, y, width, height int) {
	return w.x, w.y, 20, 20
}

func (w *RequestsWindow) SetKeybindings(ui UI) error {
	if err := ui.NewKeyBinding(w.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
		v.MoveCursor(0, 1, false)
		return nil
	}); err != nil {
		errors.Join(err)
	}

	if err := ui.NewKeyBinding(w.Name(), 'k', func(g *gocui.Gui, v *gocui.View) error {
		v.MoveCursor(0, -1, false)
		return nil
	}); err != nil {
		errors.Join(err)
	}

	if err := ui.NewKeyBinding("", '1', func(g *gocui.Gui, v *gocui.View) error {
		_, err := g.SetCurrentView(w.Name())
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
		_, err := g.SetCurrentView(ui.Views.RequestDetailsWindow.Window.Name())
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil

	// if err := ui.NewKeyBinding(w.Name(), 'n', func(g *gocui.Gui, v *gocui.View) error {
	// 	window, err := views.NewAddNewRequestView(g)
	// 	window.OnAddRequest = func(r types.Request) {
	// 		w.Requests = append(w.Requests, r)
	// 	}
	// 	return err
	// }); err != nil {
	// 	log.Panicln(err)
	// }
}
