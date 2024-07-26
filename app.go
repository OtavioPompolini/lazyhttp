package main

import (
	"errors"

	"github.com/OtavioPompolini/project-postman/ui"
	"github.com/jroimartin/gocui"
)

type App struct {
	Gui   *ui.UI
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
	app.Gui.SetFgColor(gocui.ColorGreen)
	app.Gui.SetSelectedFgColor(gocui.ColorYellow)
	app.Gui.SetCursor(true)

	if err := app.Gui.StartViews(); err != nil {
		return nil, err
	}

	if err := app.Gui.SetKeybindings(); err != nil {
		return nil, errors.Join(err)
	}

	if err := app.Gui.SetGlobalKeybindings(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run() error {
	defer app.Gui.Close()
	if err := app.Gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

// DEBUGGER
// var debuggerContent = "PUDIM"
//
// type DebuggerView struct {
// 	Title         string
// 	Name          string
// 	X, Y          int
// 	Width, Height int
// }
//
// func (*DebuggerView) NewDebuggerView() {
//
// }
//
// func NewDebuggerView(g *gocui.Gui) *DebuggerView {
// 	x, y := g.Size()
// 	return &DebuggerView{
// 		Name:   "debugger",
// 		X:      0,
// 		Y:      y - 5,
// 		Width:  x - 1,
// 		Height: 4,
// 		Title:  "Debugger",
// 	}
// }
//
// func (w *DebuggerView) Layout(g *gocui.Gui) error {
// 	v, err := g.SetView(w.Name, w.X, w.Y, w.X+w.Width, w.Y+w.Height)
// 	if err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
// 	}
//
// 	v.Clear()
// 	fmt.Fprint(v, debuggerContent)
// 	// fmt.Fprintln(v, Requests)
// 	return nil
// }
