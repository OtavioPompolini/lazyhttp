package ui

import (
	"fmt"
	"io"
	"math"

	"github.com/awesome-gocui/gocui"
)

type UI struct {
	g             *gocui.Gui
	windows       map[string]*Window
	currentWindow *Window
}

func NewUI() (*UI, error) {
	g, err := gocui.NewGui(gocui.OutputTrue, true)
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

func (ui *UI) Mouse(b bool) {
	ui.g.Mouse = b
}

func (ui *UI) SetGlobalKeybindings() error {
	for _, win := range ui.windows {
		err := win.Window.SetKeybindings(ui, win)
		if err != nil {
			return err
		}
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

func (ui *UI) SelectWindowByName(wName string) (*Window, error) {
	w, err := ui.GetWindow(wName)
	if err != nil {
		return nil, err
	}

	ui.SelectWindow(w)

	return w, nil
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
						return fmt.Errorf("error rendering window", err)
					}
				}
			}
			return nil
		}()
	})
}

func (ui *UI) renderWindow(window *Window) error {
	x, y, w, h := ui.calculateWindowSize(window)
	name := window.Window.Name()
	v, err := ui.g.SetView(name, x, y, w, h, 0)
	if err != nil {
		if err == gocui.ErrUnknownView {
			window.setView(v)
			window.Window.Setup(ui, window)
			return nil
		}
		return fmt.Errorf("Error rendering window=%s : %w", name, err)
	}
	window.Window.Update(*ui, *window)
	return nil
}

func (ui *UI) calculateWindowSize(window *Window) (x, y, w, h int) {
	wp := window.Window.Size()
	a, b := ui.Size()

	x = ui.windowSize(wp.x.position, wp.x.coord, a, 0)
	y = ui.windowSize(wp.y.position, wp.y.coord, b, 0)
	w = ui.windowSize(wp.w.position, wp.w.coord+wp.x.coord, a, 1)

	//FIX THIS
	if wp.h.position == FIXED {
		h = ui.windowSize(wp.h.position, wp.h.coord+y, b, 1)
	} else {
		h = ui.windowSize(wp.h.position, wp.h.coord+wp.y.coord, b, 1)
	}

	return x, y, w, h
}

func (ui *UI) windowSize(position position, coord int, windowSize int, off int) int {
	switch position {
	case FIXED:
		return coord
	case RELATIVE:
		v := (float64(coord) * float64(windowSize)) / 100

		//FIX THIS TOO
		if v != math.Trunc(v) {
			if off == 0 {
				v = math.Ceil(v)
			} else {
				v = math.Floor(v)
			}
		} else if off == 1 {
			v -= 1
		}

		return int(v)
	default:
		return coord
	}
}

func (ui *UI) DeleteWindow(window *Window) error {
	window.isActive = false
	if err := ui.g.DeleteView(window.Window.Name()); err != nil {
		return err
	}

	return nil
}

func (ui *UI) DeleteWindowByName(wName string) error {
	w, err := ui.GetWindow(wName)
	if err != nil {
		return err
	}

	ui.DeleteWindow(w)

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

func (ui *UI) CursorVisible(b bool) {
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

func (ui *UI) Update(f func()) {
	ui.g.Update(func(gui *gocui.Gui) error {
		f()
		return nil
	})
}

func (ui *UI) SetDefaultOutput(wName string, f func(w io.Writer)) {
	win, _ := ui.GetWindow(wName)
	f(win.view)
}
