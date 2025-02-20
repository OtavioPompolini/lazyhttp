package app

import (
	"strconv"

	"github.com/OtavioPompolini/project-postman/request"
	"github.com/OtavioPompolini/project-postman/ui"
	"github.com/jroimartin/gocui"
)

type RequestsWindow struct {
	name     string
	x, y     int
	requests *[]request.Request
	loadRequests func() error
	// OnOpenAddNewRequest func() error
}

func NewRequestsWindow(GUI *ui.UI, requests *[]request.Request) *ui.Window {
	return ui.NewWindow(
		&RequestsWindow{
			name: "RequestsWindow",
			x:    0,
			y:    0,
			requests: requests,
		})
}

func (w RequestsWindow) Name() string {
	return w.name
}

func (w *RequestsWindow) Setup(v *ui.Window) {
	v.SetSelectedBgColor(gocui.ColorRed)
	v.SetHightlight(true)
	// w.loadRequests()
}

func (w *RequestsWindow) Update(v *ui.Window) {
	v.ClearWindow()
	lines := []string{}

	for _, val := range *w.requests {
		lines = append(lines, strconv.FormatUint(uint64(val.Id), 10) + val.Name)
	}

	v.WriteLines(lines)
}

func (w *RequestsWindow) Size() (x, y, width, height int) {
	return w.x, w.y, 20, 20
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
