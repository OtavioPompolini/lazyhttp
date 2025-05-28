package app

import (
	"log"
	"math/rand/v2"
	"strconv"

	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type RequestsWindow struct {
	name           string
	windowPosition ui.WindowPosition

	// requestSystem *state.RequestSystem
	requestSystem *state.RequestSystem
	thisWindow    *ui.Window
}

func NewRequestsWindow(GUI *ui.UI, state *state.State) *ui.Window {
	requestsWindow := &RequestsWindow{
		requestSystem: state.RequestSystem,
		name:          "RequestsWindow",
		windowPosition: ui.NewWindowPosition(
			0, 20, 20, 40,
			ui.RELATIVE, ui.RELATIVE, ui.RELATIVE, ui.RELATIVE,
		),
	}

	windowRef := ui.NewWindow(
		requestsWindow,
		true,
	)

	requestsWindow.thisWindow = windowRef

	return windowRef
}

func (w *RequestsWindow) OnUpdateRequest() {
	w.thisWindow.ClearWindow()

	w.thisWindow.WriteLines(w.requestSystem.ListNames())

	err := w.thisWindow.SetCursor(0, w.requestSystem.CurrentPos())
	if err != nil {
		log.Panic(err)
	}
}

func (w *RequestsWindow) Name() string {
	return w.name
}

func (w *RequestsWindow) Setup(ui *ui.UI) {
	// ui.SelectWindow(w.thisWindow)
	w.thisWindow.SetTitle("Requests")
	w.thisWindow.SetSelectedBgColor(gocui.ColorRed)
	w.thisWindow.SetHightlight(true)

	w.requestSystem.SubscribeUpdateRequestEvent(w)
}

func (w *RequestsWindow) Update(ui ui.UI) {
}

func (w *RequestsWindow) Size() ui.WindowPosition {
	return w.windowPosition
}

func (w *RequestsWindow) SetKeybindings(ui *ui.UI) error {
	if err := ui.NewKeyBinding(w.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
		w.requestSystem.SelectNext()
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'k', func(g *gocui.Gui, v *gocui.View) error {
		w.navigateUp(ui)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'P', func(g *gocui.Gui, v *gocui.View) error {
		w.doRequest(ui)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'D', func(g *gocui.Gui, v *gocui.View) error {
		w.deleteRequest(ui)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'n', func(g *gocui.Gui, v *gocui.View) error {
		w.requestSystem.Create("pudim" + strconv.FormatInt(int64(rand.IntN(100)), 10))
		// win, _ := ui.GetWindow("CreateRequestWindow")
		// win.OpenWindow()
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'r', func(g *gocui.Gui, v *gocui.View) error {
		ui.SelectWindowByName("ResponseWindow")
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), gocui.KeyCtrlD, func(g *gocui.Gui, v *gocui.View) error {
		win, _ := ui.GetWindow("ResponseWindow")
		win.MoveCursorHalfWindowDown()

		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), gocui.KeyCtrlU, func(g *gocui.Gui, v *gocui.View) error {
		win, _ := ui.GetWindow("ResponseWindow")
		win.MoveCursorHalfWindowUp()

		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), '!', func(g *gocui.Gui, v *gocui.View) error {
		ui.SelectWindowByName("CollectionsWindow")
		return nil
	}); err != nil {
		return err
	}

	// TODO: BUT I STILL HAVEN'T FOUND WHAT I'M LOOKING FOR...
	// Handle change window with a "const" and not a string
	// and need to abstract gocui
	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
		// if w.stateService.RequestsStateService.SelectedRequest() == nil {
		// 	return nil
		// }

		ui.SelectWindowByName("RequestDetailsWindow")

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *RequestsWindow) OnDeselect(ui ui.UI, v ui.Window) error {
	return nil
}

func (w *RequestsWindow) OnSelect(ui ui.UI, v ui.Window) error {
	return nil
}

// Doest make sense to pass *ui.Window in this method
// func (w *RequestsWindow) ReloadContent(ui *ui.UI, v *ui.Window) {
// 	thisWindow, _ := ui.GetWindow(w.name)
// 	thisWindow.ClearWindow()
//
// 	lines := w.stateService.RequestsStateService.ListNames()
// 	cursorPosition := w.stateService.RequestsStateService.Index()
//
// 	thisWindow.WriteLines(lines)
//
// 	err := thisWindow.SetCursor(0, cursorPosition)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// }

// =============== ACTIONS ======================

func (rw *RequestsWindow) doRequest(ui *ui.UI) {
	// if rw.stateService.state.collection.selected == nil {
	// 	rw.stateService.state.alertMessage = "Create a request first"
	// 	win, _ := ui.GetWindow("AlertWindow")
	// 	win.OpenWindow()
	// 	return
	// }
	//
	// err := rw.stateService.ExecuteRequest()
	// if err != nil {
	// 	rw.stateService.SetAlertMessage(err.Error())
	//
	// 	alertWindow, _ := ui.GetWindow("AlertWindow")
	// 	alertWindow.OpenWindow()
	// }
	//
	// win, _ := ui.GetWindow("ResponseWindow")
	// win.Window.ReloadContent(ui, win)
}

func (rw *RequestsWindow) navigateDown(ui *ui.UI) {
	// thisWindow, _ := ui.GetWindow(rw.name)
	//
	// ok := rw.stateService.SelectNext()
	// if !ok {
	// 	return
	// }
	//
	// rw.ReloadContent(ui, thisWindow)
	// win, _ := ui.GetWindow("ResponseWindow")
	// win.Window.ReloadContent(ui, win)
}

func (rw *RequestsWindow) navigateUp(ui *ui.UI) {
	// thisWindow, _ := ui.GetWindow(rw.name)
	//
	// ok := rw.stateService.SelectPrev()
	// if !ok {
	// 	return
	// }
	//
	// rw.ReloadContent(ui, thisWindow)
	// win, _ := ui.GetWindow("ResponseWindow")
	// win.Window.ReloadContent(ui, win)
}

func (rw *RequestsWindow) deleteRequest(ui *ui.UI) {
	// thisWindow, _ := ui.GetWindow(rw.name)
	//
	// rw.stateService.DeleteSelectedRequest()
	// rw.ReloadContent(ui, thisWindow)
	//
	// win, _ := ui.GetWindow("ResponseWindow")
	// win.Window.ReloadContent(ui, win)
	// win, _ = ui.GetWindow("RequestDetailsWindow")
	// win.Window.ReloadContent(ui, win)
}
