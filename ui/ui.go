package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

type UI struct {
	g     *gocui.Gui
}

func NewUI() (*UI, error) {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return nil, err
	}

	return &UI{
		g: g,
	}, nil
}

func (ui *UI) SetHightlight(e bool) {
	ui.g.Highlight = e
}

func (ui *UI) SetGlobalKeybindings() error {
	if err := ui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}

	if err := ui.g.SetKeybinding("", gocui.KeyCtrl2, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		g.SetCurrentView("requestsWindow")
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (ui *UI) SelectWindow(viewName string) error {
	if _, err := ui.g.SetCurrentView(viewName); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (ui *UI) SetFgColor(color gocui.Attribute) {
	ui.g.FgColor = color
}

func (ui *UI) SetSelectedFgColor(color gocui.Attribute) {
	ui.g.SelFgColor = color
}

func (ui *UI) SetManagerFunc(f func() error) {
	ui.g.SetManagerFunc(func(*gocui.Gui) error {
		return f()
	})
}

func (ui *UI) RenderWindow(window *Window) error {
	x, y, w, h := window.Window.Size()
	name := window.Window.Name()
	v, err := ui.g.SetView(name, x, y, w, h)
	if err != nil {
		if err == gocui.ErrUnknownView {
			window.setView(v)
			window.Window.Setup(window)
			return nil
		}
		return fmt.Errorf("Error rendering window=%s : %w", name, err)
	}
	window.Window.Update(window)
	return nil
}

func (ui *UI) Start() error {
	return ui.g.MainLoop()
}

func (ui *UI) NewKeyBinding(name string, key interface{}, callback func(g *gocui.Gui, v *gocui.View) error) error {
	if err := ui.g.SetKeybinding(name, key, gocui.ModNone, callback); err != nil {
		return err
	}
	return nil
}

func (ui *UI) Close() {
	ui.g.Close()
}

func (ui *UI) SetCursor(b bool) {
	ui.g.Cursor = b
}

func (ui *UI) ActiveViewName() string {
	return ui.g.CurrentView().Name()
}
