package app

import (
	"github.com/OtavioPompolini/project-postman/memory"
	"github.com/OtavioPompolini/project-postman/ui"
	"github.com/jroimartin/gocui"
)

type RequestsWindow struct {
	isActive     bool
	name         string
	x, y, h, w   int
	memory       *memory.Memory
	loadRequests func() error
	currentLine  int
	// OnOpenAddNewRequest func() error
}

func NewRequestsWindow(GUI *ui.UI, memomry *memory.Memory) *ui.Window {
	_, b := GUI.Size()
	return ui.NewWindow(
		&RequestsWindow{
			name:     "RequestsWindow",
			x:        0,
			y:        0,
			h:        b - 1,
			w:        49,
			memory:   memomry,
			isActive: true,
		},
		true,
	)
}

func (w RequestsWindow) Name() string {
	return w.name
}

func (w *RequestsWindow) Setup(ui ui.UI, v ui.Window) {
	ui.SelectWindow(&v)
	v.SetTitle(v.Window.Name())
	v.SetSelectedBgColor(gocui.ColorRed)
	v.SetHightlight(true)
}

func (w *RequestsWindow) Update(ui ui.UI, v ui.Window) {
	v.ClearWindow()

	requests := w.memory.ListRequests()

	lines := []string{}

	for _, r := range requests {
		lines = append(lines, r.Name)
	}

	v.WriteLines(lines)
}

func (w *RequestsWindow) Size() (x, y, width, height int) {
	return w.x, w.y, w.x + w.w, w.y + w.h
}

func (w *RequestsWindow) IsActive() bool {
	return w.isActive
}

func (w *RequestsWindow) SetKeybindings(ui *ui.UI) error {

	if err := ui.NewKeyBinding(w.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
		w.memory.SelectNext()
		v.MoveCursor(0, 1, false)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'k', func(g *gocui.Gui, v *gocui.View) error {
		w.memory.SelectPrev()
		v.MoveCursor(0, -1, false)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'n', func(g *gocui.Gui, v *gocui.View) error {
		win, err := ui.GetWindow("CreateRequestWindow")
		if err != nil {
			return err
		}

		win.SwitchOnOff(true)

		return nil
	}); err != nil {
		return err
	}

	//TODO: BUT I STILL HAVEN'T FOUND WHAT I'M LOOKING FOR...
	//Handle change window with a "const" and not a string
	// and need to abstract gocui
	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
		if w.memory.IsEmpty() {
			return nil
		}

		toWindow, err := ui.GetWindow("RequestDetailsWindow")
		if err != nil {
			return err
		}
		_, err = ui.SelectWindow(toWindow)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *RequestsWindow) OnDeselect(ui ui.UI, v ui.Window) error {
	return nil
}

func (w *RequestsWindow) OnSelect(ui ui.UI, v ui.Window) error {
	return nil
}

// func (w *RequestsWindow) SetKeybindings(ui ui.UI) error {
//
// 	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
// 		// _, err := g.SetCurrentView(ui.Views.RequestDetailsWindow.Window.Name())
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		return nil
// 	}); err != nil {
// 		return err
// 	}
//
// 	return nil
//
// 	// if err := ui.NewKeyBinding(w.Name(), 'n', func(g *gocui.Gui, v *gocui.View) error {
// 	// 	window, err := views.NewAddNewRequestView(g)
// 	// 	window.OnAddRequest = func(r types.Request) {
// 	// 		w.Requests = append(w.Requests, r)
// 	// 	}
// 	// 	return err
// 	// }); err != nil {
// 	// 	log.Panicln(err)
// 	// }
// }
