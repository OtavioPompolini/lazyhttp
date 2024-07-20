package main

import (
	"fmt"
	"log"

	"github.com/OtavioPompolini/project-postman/ui"
	"github.com/OtavioPompolini/project-postman/views"
	"github.com/jroimartin/gocui"
)

type Views struct {
	RequestsView *views.RequestsView
}


type App struct {
	Gui *ui.UI
	Views *Views
}

func NewApp() (*App, error) {

	userInteface, err := ui.NewUI();
	if  err != nil {
		log.Panicln(err)
	}

	app := &App{
		Gui: userInteface,
		Views: &Views{
			RequestsView: views.NewRequestsView(),
		},
	}

	app.Gui.SetHightlight(true)
	app.Gui.SetFgColor(gocui.ColorBlue)

	// debuggerView := NewDebuggerView()
	// requestsView := NewRequestsView()
	// RequestDetailView := NewRequestDetailView(g)

	app.Gui.SetManager()

	// g.SetManagerFunc()
	// g.SetManager(requestsView, debuggerView, RequestDetailView)

	// if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
	// 	return gocui.ErrQuit
	// }); err != nil {
	// 	log.Panicln(err)
	// }

	// if err := g.SetKeybinding(arw.Name, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
	// 	read_line := v.Buffer()
	// 	read_line = strings.TrimSuffix(read_line, "\n")
	// 	arw.OnAddRequest(types.NewRequest(read_line))
	// 	g.DeleteView("addNewRequest")
	// 	return nil
	// }); err != nil {
	// 	log.Panicln(err)
	// }

	// if err := g.SetKeybinding("addNewRequest", gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
	// 	return ChangeView(g, LastViewId)
	// }); err != nil {
	// 	log.Panicln(err)
	// }

	return app, nil
}

func (app *App) StartViews() error {
	views := []ui.View{}
	app.Gui.SetManager(views...)
	return nil
}

func (app *App) Run() error {
	defer app.Gui.Close()
	if err := app.Gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
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
		Name:   "debugger",
		X:      0,
		Y:      y - 5,
		Width:  x - 1,
		Height: 4,
		Title:  "Debugger",
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
	// fmt.Fprintln(v, Requests)
	return nil
}

