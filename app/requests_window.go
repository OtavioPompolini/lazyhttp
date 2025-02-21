package app

import (
	"sort"
	"strconv"

	"github.com/OtavioPompolini/project-postman/request"
	"github.com/OtavioPompolini/project-postman/ui"
	"github.com/jroimartin/gocui"
)

type RequestsWindow struct {
	isActive     bool
	name         string
	x, y         int
	requests     *map[int64]request.Request
	loadRequests func() error
	currentLine  int
	// OnOpenAddNewRequest func() error
}

func NewRequestsWindow(GUI *ui.UI, requests *map[int64]request.Request) *ui.Window {
	return ui.NewWindow(
		&RequestsWindow{
			name:     "RequestsWindow",
			x:        0,
			y:        0,
			requests: requests,
			isActive: true,
		})
}

func (w RequestsWindow) Name() string {
	return w.name
}

func (w *RequestsWindow) Setup(v ui.Window) {
	v.SetTitle(v.Window.Name())
	v.SetSelectedBgColor(gocui.ColorRed)
	v.SetHightlight(true)
	// w.loadRequests()
}

func (w *RequestsWindow) Update(v ui.Window) {
	v.ClearWindow()
	lines := []string{}

	for _, val := range *w.requests {
		lines = append(lines, strconv.FormatInt(val.Id, 10)+val.Name)
	}

	sort.Strings(lines)
	v.WriteLines(lines)
}

func (w *RequestsWindow) Size() (x, y, width, height int) {
	return w.x, w.y, 20, 20
}

func (w *RequestsWindow) IsActive() bool {
	return w.isActive
}

func (w *RequestsWindow) SetKeybindings(ui ui.UI) error {

	if err := ui.NewKeyBinding(w.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
		v.MoveCursor(0, 1, false)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'k', func(g *gocui.Gui, v *gocui.View) error {
		v.MoveCursor(0, -1, false)
		return nil
	}); err != nil {
		return err
	}

	//TODO: BUT I STILL HAVEN'T FOUND WHAT I'M LOOKING FOR...
	//Handle change window with a "const" and not a string
	// and need to abstract gocui
	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
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

func (w *RequestsWindow) OnDeselect() error {
	return nil
}

func (w *RequestsWindow) OnSelect() error {
	return nil
}

// func (w *RequestsWindow) SetKeybindings(ui ui.UI) error {
// 	if err := ui.NewKeyBinding(w.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
// 		v.MoveCursor(0, 1, false)
// 		return nil
// 	}); err != nil {
// 		errors.Join(err)
// 	}
//
// 	if err := ui.NewKeyBinding(w.Name(), 'k', func(g *gocui.Gui, v *gocui.View) error {
// 		v.MoveCursor(0, -1, false)
// 		return nil
// 	}); err != nil {
// 		errors.Join(err)
// 	}
//
// 	if err := ui.NewKeyBinding("", '1', func(g *gocui.Gui, v *gocui.View) error {
// 		_, err := g.SetCurrentView(w.Name())
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}); err != nil {
// 		return err
// 	}
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
