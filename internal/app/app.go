package app

import (
	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/memory"
	"github.com/OtavioPompolini/project-postman/internal/ui"
)
type App struct {
	GUI *ui.UI
	debuggerMode  bool
	memoryHandler *memory.Memory
}

func NewApp() (*App, error) {
	userInteface, err := ui.NewUI()
	if err != nil {
		return nil, err
	}

	memory, err := memory.InitMemory()
	if err != nil {
		return nil, err
	}

	app := &App{
		GUI:           userInteface,
		memoryHandler: memory,
		debuggerMode:  true, // TODO: run argument --debug=true
	}

	app.GUI.AddWindow(NewRequestsWindow(userInteface, app.memoryHandler))
	app.GUI.AddWindow(NewRequestDetailsWindow(userInteface, app.memoryHandler))
	app.GUI.AddWindow(NewCreateRequestWindow(userInteface, app.memoryHandler))
	app.GUI.AddWindow(NewResponseWindow(userInteface, app.memoryHandler))

	app.GUI.StartUI()

	app.GUI.SetHightlight(true)
	app.GUI.SetFgColor(gocui.ColorGreen)
	app.GUI.SetSelectedFgColor(gocui.ColorYellow)

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

// DEBUGGER
// var debuggerContent = "PUDIM"
//
// type DebuggerView struct {
// 	Title string
// 	x, y  int
// 	w, h  int
// }
//
// func NewDebuggerView(GUI *ui.UI) *ui.Window {
// 	x, y := GUI.Size()
//
// 	return ui.NewWindow(
// 		&DebuggerView{
// 			Title: "debugger",
// 			x:     0,
// 			y:     y - 5,
// 			h:     4,
// 			w:     x - 1,
// 		})
// }
//
// func (w *DebuggerView) Setup(v *ui.Window) {
// 	v.SetSelectedBgColor(gocui.ColorRed)
// 	v.SetHightlight(true)
// }
//
// func (w *DebuggerView) Update(v *ui.Window) {
// 	v.ClearWindow()
// 	v.WriteLn(debuggerContent)
// }
//
// func (w *DebuggerView) Name() string {
// 	return "DebuggerView"
// }
//
// func (w *DebuggerView) Size() (x, y, wid, hei int) {
// 	return w.x, w.y, w.x + w.w, w.y + w.h
// }
//
// func (w *DebuggerView) OnDeselect() error {
// 	return nil
// }
//
// func (w *DebuggerView) OnSelect() error {
// 	return nil
// }
