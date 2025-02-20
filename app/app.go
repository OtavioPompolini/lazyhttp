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
	CreateRequestWindow  *ui.Window
	DebuggerView         *ui.Window
}

type App struct {
	GUI          *ui.UI
	db           request.Adapter
	Windows      *Windows
	debuggerMode bool
	requests     *[]request.Request
}

func NewApp() (*App, error) {

	userInteface, err := ui.NewUI()
	if err != nil {
		return nil, err
	}

	db, err := request.InitDatabase()
	if err != nil {
		return nil, err
	}

	requests := db.GetRequests()

	app := &App{
		GUI:      userInteface,
		db:       db,
		requests: &requests, //RENAME MEMORY AND CREATE STRUCT
		Windows: &Windows{
			RequestsWindow:       NewRequestsWindow(userInteface, &requests),
			RequestDetailsWindow: NewRequestDetailsWindow(userInteface),
			CreateRequestWindow: NewCreateRequestWindow(userInteface, func(reqName string) {
				db.CreateRequest(reqName)
			}),
			DebuggerView: NewDebuggerView(userInteface),
		},
		debuggerMode: true, //TODO: run argument --debug=true
	}

	app.GUI.SetManagerFunc(app.layout)
	app.GUI.SetHightlight(true)
	app.GUI.SetFgColor(gocui.ColorGreen)
	app.GUI.SetSelectedFgColor(gocui.ColorYellow)

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

	if app.debuggerMode == true {
		if err := app.GUI.RenderWindow(app.Windows.DebuggerView); err != nil {
			return err
		}
	}

	// if app.isCreatingRequest == true {
	// if err := app.GUI.RenderWindow(app.Windows.CreateRequestWindow); err != nil {
	// 	return err
	// }
	// }

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

	if err := app.GUI.NewKeyBinding(app.Windows.RequestsWindow.Window.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
		_, err := app.GUI.SelectWindow(app.Windows.RequestDetailsWindow)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		errors.Join(err)
	}

	if err := app.GUI.NewKeyBinding(app.Windows.RequestDetailsWindow.Window.Name(), gocui.KeyEsc, func(g *gocui.Gui, v *gocui.View) error {
		_, err := app.GUI.SelectWindow(app.Windows.RequestsWindow)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		errors.Join(err)
	}

	if err := app.GUI.NewKeyBinding("", '1', func(g *gocui.Gui, v *gocui.View) error {
		_, err := app.GUI.SelectWindow(app.Windows.RequestsWindow)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := app.GUI.NewKeyBinding("", '2', func(g *gocui.Gui, v *gocui.View) error {
		_, err := app.GUI.SelectWindow(app.Windows.RequestsWindow)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := app.GUI.NewKeyBinding("", '3', func(g *gocui.Gui, v *gocui.View) error {
		saved := app.db.CreateRequest("PUDIM")
		*app.requests = append(*app.requests, saved)
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// DEBUGGER
var debuggerContent = "PUDIM"

type DebuggerView struct {
	Title string
	x, y  int
	w, h  int
}

func NewDebuggerView(GUI *ui.UI) *ui.Window {
	x, y := GUI.Size()

	return ui.NewWindow(
		&DebuggerView{
			Title: "debugger",
			x:     0,
			y:     y - 5,
			h:     4,
			w:     x - 1,
		})
}

func (w *DebuggerView) Setup(v *ui.Window) {
	v.SetSelectedBgColor(gocui.ColorRed)
	v.SetHightlight(true)
}

func (w *DebuggerView) Update(v *ui.Window) {
	v.ClearWindow()
	v.WriteLn(debuggerContent)
}

func (w *DebuggerView) Name() string {
	return "DebuggerView"
}

func (w *DebuggerView) Size() (x, y, wid, hei int) {
	return w.x, w.y, w.x + w.w, w.y + w.h
}

func (w *DebuggerView) OnDeselect() error {
	return nil
}

func (w *DebuggerView) OnSelect() error {
	return nil
}
