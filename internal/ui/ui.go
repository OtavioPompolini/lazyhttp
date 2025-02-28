package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type UI struct {
	g             *gocui.Gui
	windows       map[string]*Window
	currentWindow *Window
}

func NewUI() (*UI, error) {
	g, err := gocui.NewGui(gocui.Output256, true)
	if err != nil {
		return nil, err
	}

	return &UI{
		g:       g,
		windows: map[string]*Window{},
	}, nil
}

func (ui *UI) AddWindow(window *Window) {
	ui.windows[window.Window.Name()] = window
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

// This method doest need window return
// Isnt better to pass window name string and get the actual Window inside here?
func (ui *UI) SelectWindow(window *Window) (*Window, error) {
	if ui.currentWindow != nil {
		ui.currentWindow.Window.OnDeselect(*ui, *ui.currentWindow)
	}

	if _, err := ui.g.SetCurrentView(window.Window.Name()); err != nil {
		return nil, err
	}

	window.Window.OnSelect(*ui, *window)

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
	ui.g.InputEsc = true
	ui.g.SetManagerFunc(func(*gocui.Gui) error {
		return func() error {
			for _, win := range ui.windows {
				if win.IsActive() {
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
	v, err := ui.g.SetView(name, x, y, w, h, 1)
	if err != nil {
		if err == gocui.ErrUnknownView {
			window.setView(v)
			window.Window.Setup(*ui, *window)
			return nil
		}
		return fmt.Errorf("Error rendering window=%s : %w", name, err)
	}
	window.Window.Update(*ui, *window)
	return nil
}

func (ui *UI) DeleteWindow(window *Window) error {
	window.isActive = false
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

func (ui *UI) GetWindow(wName string) (*Window, error) {
	w, ok := ui.windows[wName]
	if ok {
		return w, nil
	}

	return nil, fmt.Errorf("View %s does not exists", wName)
}

func (ui *UI) Size() (int, int) {
	return ui.g.Size()
}
