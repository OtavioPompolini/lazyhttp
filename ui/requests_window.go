package ui

import (
	"github.com/jroimartin/gocui"
)

type RequestsWindow struct{
	OnClick func() error
}

func (RequestsWindow) Name() string {
	return "requestsWindow"
}

func NewRequestsWindow(ui *UI) *RequestsWindow {
	requestsWindow := &RequestsWindow{}
	return requestsWindow
}

func (w *RequestsWindow) Setup(v *View) {
	v.SetSelectedBgColor(gocui.ColorRed)
	v.SetHightlight(true)
	for i:=0;i<10;i++ {
		v.WriteLn("PUDIM")
	}
}

func (w *RequestsWindow) Update(v *View) {
}

func (w *RequestsWindow) Size() (x, y, width, height int) {
	return 0, 0, 20, 20
}

func (w *RequestsWindow) SetKeybindings(ui *UI) error {
	if err := ui.NewKeyBinding(w.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
		v.MoveCursor(0, 1, false)
		return nil
	}); err != nil {
		return err
	}
	return nil
	//
	// if err := ui.NewKeyBinding(w.Name(), 'k', func(g *gocui.Gui, v *gocui.View) error {
	// 	v.MoveCursor(0, -1, false)
	// 	_, y := v.Cursor()
	// 	SelectedRequest = y
	// 	return nil
	// }); err != nil {
	// 	log.Panicln(err)
	// }
	//
	// if err := ui.NewKeyBinding(w.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
	// 	ChangeView(g, RequestDetailViewName)
	// 	return nil
	// }); err != nil {
	// 	log.Panicln(err)
	// }
	//
	// if err := ui.NewKeyBinding("", '1', func(g *gocui.Gui, v *gocui.View) error {
	// 	ChangeView(g, w.Name)
	// 	return nil
	// }); err != nil {
	// 	log.Panicln(err)
	// }
	//
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
