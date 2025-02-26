package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type UI struct {
	g             *gocui.Gui
	windows       []*Window
	currentWindow *Window
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

func (ui *UI) AddWindow(window *Window) {
	ui.windows = append(ui.windows, window)
}

func (ui *UI) SetHightlight(e bool) {
	ui.g.Highlight = e
}

func (ui *UI) SetGlobalKeybindings() error {
	for _, win := range ui.windows {
		win.Window.SetKeybindings(ui)
	}

	if err := ui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}

	return nil
}

func (ui *UI) SelectWindow(window *Window) (*Window, error) {
	if ui.currentWindow != nil {
		ui.currentWindow.Window.OnDeselect()
	}

	if _, err := ui.g.SetCurrentView(window.Window.Name()); err != nil {
		return nil, err
	}

	window.Window.OnSelect()

	ui.currentWindow = window
	return window, nil
}

func (ui *UI) SetFgColor(color gocui.Attribute) {
	ui.g.FgColor = color
}

func (ui *UI) SetSelectedFgColor(color gocui.Attribute) {
	ui.g.SelFgColor = color
}

func (ui *UI) StartUI() {
	ui.g.SetManagerFunc(func(*gocui.Gui) error {
		return func() error {
			for _, win := range ui.windows {
				if win.Window.IsActive() {
					if err := ui.renderWindow(win); err != nil {
						return fmt.Errorf("error rendering window")
					}
				}
			}
			return nil
		}()
	})
}

func (ui *UI) renderWindow(window *Window) error {
	x, y, w, h := window.Window.Size()
	name := window.Window.Name()
	v, err := ui.g.SetView(name, x, y, w, h)
	if err != nil {
		if err == gocui.ErrUnknownView {
			window.setView(v)
			window.Window.Setup(*window)
			return nil
		}
		return fmt.Errorf("Error rendering window=%s : %w", name, err)
	}
	window.Window.Update(*window)
	return nil
}

func (ui *UI) DeleteWindow(window *Window) error {
	if err := ui.g.DeleteView(window.Window.Name()); err != nil {
		return err
	}

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

func (ui *UI) GetWindow(vName string) (*Window, error) {
	for _, win := range ui.windows {
		if win.Window.Name() == vName {
			return win, nil
		}
	}

	return nil, fmt.Errorf("View %s does not exists", vName)
}

func (ui *UI) Size() (int, int) {
	return ui.g.Size()
}
