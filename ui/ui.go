package ui

import (
	"github.com/jroimartin/gocui"
)

type UI struct {
	g *gocui.Gui
}

func NewUI() (*UI, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	g.Highlight = true
	g.FgColor = gocui.ColorCyan

	return &UI{
		g: g,
	}, nil
}

func (ui *UI) SetHightlight(e bool) {
	ui.g.Highlight = e
}

func (ui *UI) SetCloseKeybinding() error {
	if err := ui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	return nil
}

func (ui *UI) SelectWindow(viewName string) error {
	if _, err := ui.g.SetCurrentView(viewName); err != nil {
		return err
	}
	return nil
}

func (ui *UI) SetFgColor(color gocui.Attribute) {
	ui.g.FgColor = color
}

func (ui *UI) SetWindows(views ...gocui.Manager) {
	ui.g.SetManager(views...)
}

func (ui *UI) MainLoop() error {
	return ui.g.MainLoop()
}

func (ui *UI) NewKeyBinding(name string, key gocui.Modifier, callback func() error) {
	ui.g.SetKeybinding(name, gocui.ModNone, key, func(g *gocui.Gui, v *gocui.View) error {
		return callback()
	})
}

func (ui *UI) Close() {
	ui.g.Close()
}

func (ui *UI) ActiveViewName() string {
	return ui.g.CurrentView().Name()
}
