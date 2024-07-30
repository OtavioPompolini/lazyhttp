package app

import (
	"errors"

	"github.com/OtavioPompolini/project-postman/request"
	"github.com/OtavioPompolini/project-postman/ui"
	"github.com/jroimartin/gocui"
)

type Windows struct {
	RequestsWindow       *ui.Window
	RequestDetailsWindow *ui.Window
}

type App struct {
	GUI     *ui.UI
	adapter *request.Adapter
	Windows *Windows
}

func NewApp() (*App, error) {

	userInteface, err := ui.NewUI()
	if err != nil {
		return nil, err
	}

	app := &App{
		GUI: userInteface,
		Windows: &Windows{
			RequestsWindow:       NewRequestsWindow(userInteface),
			RequestDetailsWindow: NewRequestDetailsWindow(userInteface),
		},
	}

	app.GUI.SetManagerFunc(app.layout)
	app.GUI.SetHightlight(true)
	app.GUI.SetFgColor(gocui.ColorGreen)
	app.GUI.SetSelectedFgColor(gocui.ColorYellow)
	app.GUI.SetCursor(true)

	if err := app.SetKeybindings(); err != nil {
		return nil, errors.Join(err)
	}

	if err := app.GUI.SetGlobalKeybindings(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run() error {
	defer app.GUI.Close()
	if err := app.GUI.Start(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (app *App) layout() error {
	if err := app.GUI.RenderWindow(app.Windows.RequestsWindow); err != nil {
		return err
	}

	if err := app.GUI.RenderWindow(app.Windows.RequestDetailsWindow); err != nil {
		return err
	}

	return nil
}


// TODO: move this to other file and abstract func(g *gocui.Gui, v *gocui.View)
func (app *App) SetKeybindings() error {

	if err := app.GUI.NewKeyBinding(app.Windows.RequestsWindow.Window.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
		v.MoveCursor(0, 1, false)
		return nil
	}); err != nil {
		errors.Join(err)
	}

	if err := app.GUI.NewKeyBinding(app.Windows.RequestsWindow.Window.Name(), 'k', func(g *gocui.Gui, v *gocui.View) error {
		v.MoveCursor(0, -1, false)
		return nil
	}); err != nil {
		errors.Join(err)
	}

	if err := app.GUI.NewKeyBinding("", '1', func(g *gocui.Gui, v *gocui.View) error {
		_, err := g.SetCurrentView(app.Windows.RequestsWindow.Window.Name())
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
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
