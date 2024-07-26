package ui

import (
	"errors"
	"log"

	"github.com/jroimartin/gocui"
)

type Views struct {
	RequestsWindow       *View
	RequestDetailsWindow *View
}

type UI struct {
	g     *gocui.Gui
	Views *Views
}

func NewUI() (*UI, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	return &UI{
		g: g,
	}, nil
}

func (ui *UI) StartViews() error {
	ui.Views = &Views{
		RequestsWindow:       NewRequestsWindow(ui),
		RequestDetailsWindow: NewRequestDetailsWindow(ui),
	}

	ui.SetWindows()
	return nil
}

func (ui *UI) SetHightlight(e bool) {
	ui.g.Highlight = e
}

func (ui *UI) SetKeybindings() error {
	if err := ui.Views.RequestsWindow.Window.SetKeybindings(*ui); err != nil {
		return errors.Join(err)
	}
	return nil
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

func (ui *UI) SetWindows() {
	ui.g.SetManager(
		ui.Views.RequestsWindow,
		ui.Views.RequestDetailsWindow,
	)
}

func (ui *UI) MainLoop() error {
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
