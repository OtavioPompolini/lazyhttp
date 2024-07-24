package main

import (
	"fmt"

	"github.com/OtavioPompolini/project-postman/ui"
	"github.com/OtavioPompolini/project-postman/views"
	"github.com/jroimartin/gocui"
)

type Views struct {
	RequestsWindow *ui.View
}

type App struct {
	Gui   *ui.UI
	Views *Views
}

func NewApp() (*App, error) {

	userInteface, err := ui.NewUI()
	if err != nil {
		return nil, err
	}

	app := &App{
		Gui: userInteface,
	}

	app.Gui.SetHightlight(true)
	app.Gui.SetFgColor(gocui.ColorBlue)

	if err := app.StartViews(); err != nil {
		return nil, err
	}

	if err := app.Gui.SetCloseKeybinding(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) StartViews() error {
	view, err := ui.NewView(
		&views.RequestsWindow{},
		"requests",
	)
	if err != nil {
		return err
	}
	app.Views = &Views{
		RequestsWindow: view,
	}

	app.Gui.SetWindows(
		app.Views.RequestsWindow,
	)
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
