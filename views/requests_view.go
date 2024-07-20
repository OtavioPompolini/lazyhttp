package views

import (
	"github.com/OtavioPompolini/project-postman/ui"
)

type RequestsView = ui.View

func NewRequestsView() *RequestsView {
	return ui.NewView("requests")
}

// func (w *RequestsView) Layout(g *gocui.Gui) error {
// 	v, err := g.SetView(w.Name, w.X, w.Y, w.X+w.Width, w.Y+w.Height)
// 	if err != nil && err != gocui.ErrUnknownView {
// 		return err
// 	}
//
// 	if err == gocui.ErrUnknownView {
// 		_, err := g.SetCurrentView(w.Name)
//
// 		v.Title = w.Title
// 		v.SelBgColor = gocui.ColorRed
// 		v.Highlight = true
//
// 		w.setKeybindings(g)
//
// 		return err
// 	}
//
// 	v.Clear()
//
// 	for i, request := range w.Requests {
// 		if i == len(w.Requests)-1 {
// 			fmt.Fprint(v, request.Name)
// 			continue
// 		}
// 		fmt.Fprintln(v, request.Name)
// 	}
//
// 	return nil
// }
//
// func (w *RequestsView) setKeybindings(g *gocui.Gui) {
// 	if err := g.SetKeybinding(RequestViewName, 'j', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
// 		v.MoveCursor(0, 1, false)
// 		_, y := v.Cursor()
// 		SelectedRequest = y
// 		// debuggerContent = fmt.Sprintln(v.Cursor())
// 		return nil
// 	}); err != nil {
// 		log.Panicln(err)
// 	}
//
// 	if err := g.SetKeybinding(RequestViewName, 'k', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
// 		v.MoveCursor(0, -1, false)
// 		_, y := v.Cursor()
// 		SelectedRequest = y
// 		return nil
// 	}); err != nil {
// 		log.Panicln(err)
// 	}
//
// 	if err := g.SetKeybinding(RequestViewName, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
// 		ChangeView(g, RequestDetailViewName)
// 		return nil
// 	}); err != nil {
// 		log.Panicln(err)
// 	}
//
// 	if err := g.SetKeybinding("", '1', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
// 		ChangeView(g, w.Name)
// 		return nil
// 	}); err != nil {
// 		log.Panicln(err)
// 	}
//
// 	if err := g.SetKeybinding(RequestViewName, 'n', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
// 		window, err := views.NewAddNewRequestView(g)
// 		window.OnAddRequest = func(r types.Request) {
// 			w.Requests = append(w.Requests, r)
// 		}
// 		return err
// 	}); err != nil {
// 		log.Panicln(err)
// 	}
// }
