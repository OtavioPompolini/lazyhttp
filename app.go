package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/OtavioPompolini/project-postman/views"
	"github.com/OtavioPompolini/project-postman/types"
	"github.com/jroimartin/gocui"
)

var (
	Requests = [](*types.Request){
		types.NewRequest("pudim"),
		types.NewRequest("pudim"),
		types.NewRequest("pudim"),
		types.NewRequest("pudim"),
		types.NewRequest("pudim"),
	}
	SelectedRequest = 0
)

// TODO: Move this outside

const RequestViewName = "requests"

type RequestsView struct {
	Title         string
	Name          string
	X, Y          int
	Width, Height int
	SelectKey     gocui.Key
}

func NewRequestsView() *RequestsView {
	//TODO: Make this responsible
	return &RequestsView{
		Name:      RequestViewName,
		X:         0,
		Y:         0,
		Width:     14,
		Height:    20,
		Title:     "Requests",
	}
}

func (w *RequestsView) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.Name, w.X, w.Y, w.X+w.Width, w.Y+w.Height)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if err == gocui.ErrUnknownView {
		_, err := g.SetCurrentView(w.Name)

		v.Title = w.Title
		v.SelBgColor = gocui.ColorRed
		v.Highlight = true

		w.setKeybindings(g);

		return err
	}

	v.Clear()

	for _, request := range Requests {
		fmt.Fprintln(v, request.Name)
	}

	// if err := g.SetKeybinding("", '1', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
	// 	_, err = g.SetCurrentView(w.Name)
	// 	return err
	// }); err != nil {
	// 	log.Panicln(err)
	// }

	return nil
}

func (w *RequestsView) setKeybindings(g *gocui.Gui) {
	if err := g.SetKeybinding(RequestViewName, 'j', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		v.MoveCursor(0, 1, false)
		// debuggerContent = fmt.Sprintln(v.Cursor())
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(RequestViewName, 'k', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		v.MoveCursor(0, -1, false)
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(RequestViewName, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, y := v.Cursor()
		SelectedRequest = y
		g.SetCurrentView(RequestDetailViewName)

		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", '1', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		g.SetCurrentView(w.Name)
		return nil
	}); err != nil {
		log.Panicln(err)
	}
}

const RequestDetailViewName = "request_detail"

type RequestDetailView struct {
	Title         string
	Name          string
	X, Y          int
	Width, Height int
}

func NewRequestDetailView(g *gocui.Gui) *RequestDetailView {
	return &RequestDetailView{
		Name:      RequestViewName,
		X:         0,
		Y:         0,
		Width:     14,
		Height:    20,
		Title:     "Requests",
	}
}

func (w *RequestDetailView) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.Name, w.X, w.Y, w.X+w.Width, w.Y+w.Height)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if err == gocui.ErrUnknownView {
		_, err := g.SetCurrentView(w.Name)

		v.Title = w.Title
		v.SelBgColor = gocui.ColorRed
		v.Highlight = true

		// w.setKeybindings(g);

		return err
	}

	v.Clear()

	for _, request := range Requests {
		fmt.Fprintln(v, request.Name)
	}

	// if err := g.SetKeybinding("", '1', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
	// 	_, err = g.SetCurrentView(w.Name)
	// 	return err
	// }); err != nil {
	// 	log.Panicln(err)
	// }

	return nil
}

// func SelectRequestsWindow(g *gocui.Gui) (*gocui.View, error) {
// 	v, err := g.SetCurrentView(RequestViewName)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	v.SetCursor(0, 0)
// 	return v, nil
// }

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	g.Highlight = true
	g.SelFgColor = gocui.ColorBlack

	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	debuggerView := NewDebuggerView(g)
	requestsView := NewRequestsView()
	RequestDetailView := NewRequestDetailView(g)
	// view2 := newWindow("env", 0, 3, 14, 2, "marcos viado", gocui.KeyF2)

	g.SetManager(requestsView, debuggerView, RequestDetailView)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("requests", 'n', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		err := views.OpenAddNewRequestView(g)
		return err
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(views.AddNewRequestName, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		read_line := v.Buffer()
		read_line = strings.TrimSuffix(read_line, "\n")
		Requests = append(Requests, types.NewRequest(read_line))
		g.DeleteView("addNewRequest")
		g.SetCurrentView("requests")
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("addNewRequest", gocui.KeyF1, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		g.SetCurrentView("requests")
		g.DeleteView("addNewRequest")
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}


// DEBUGGER
var debuggerContent = "PUDIM"

type DebuggerView struct {
	Title         string
	Name          string
	X, Y          int
	Width, Height int
}

func (*DebuggerView) NewDebuggerView() {

}

func NewDebuggerView(g *gocui.Gui) *DebuggerView {
	x, y := g.Size()
	return &DebuggerView{
		Name:      "debugger",
		X:         0,
		Y:         y-5,
		Width:     x-1,
		Height:    4,
		Title:     "Debugger",
	}
}

func (w *DebuggerView) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.Name, w.X, w.Y, w.X+w.Width, w.Y+w.Height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	v.Clear()
	fmt.Fprint(v, debuggerContent)
	return nil
}

