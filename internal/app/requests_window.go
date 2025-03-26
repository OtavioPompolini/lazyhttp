package app

import (
	"log"

	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/ui"
)

// Why not to save a ref to ui.GUI and ui.Window on every Window implementation?????
// Since every F method needs both params...
type RequestsWindow struct {
	isActive     bool
	name         string
	x, y, h, w   int
	stateService StateService
	loadRequests func() error
}

func NewRequestsWindow(GUI *ui.UI, stateService StateService) *ui.Window {
	a, b := GUI.Size()
	return ui.NewWindow(
		&RequestsWindow{
			name:         "RequestsWindow",
			x:            0,
			y:            0,
			h:            b - 1,
			w:            a * 20 / 100,
			stateService: stateService,
			isActive:     true,
		},
		true,
	)
}

func (w RequestsWindow) Name() string {
	return w.name
}

func (w *RequestsWindow) Setup(ui *ui.UI, v *ui.Window) {
	ui.SelectWindow(v)
	v.SetTitle("Requests:")
	v.SetSelectedBgColor(gocui.ColorRed)
	v.SetHightlight(true)
	w.ReloadContent(ui, v)
}

func (w *RequestsWindow) Update(ui ui.UI, v ui.Window) {
}

func (w *RequestsWindow) Size() (x, y, width, height int) {
	return w.x, w.y, w.x + w.w, w.y + w.h
}

func (w *RequestsWindow) IsActive() bool {
	return w.isActive
}

func (w *RequestsWindow) SetKeybindings(ui *ui.UI, win *ui.Window) error {
	if err := ui.NewKeyBinding(w.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
		w.navigateDown(ui)
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
		win, _ := ui.GetWindow("CreateRequestWindow")
		win.OpenWindow()
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

	// TODO: BUT I STILL HAVEN'T FOUND WHAT I'M LOOKING FOR...
	// Handle change window with a "const" and not a string
	// and need to abstract gocui
	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
		if w.stateService.state.collection.selected == nil {
			return nil
		}

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
func (w *RequestsWindow) ReloadContent(ui *ui.UI, v *ui.Window) {
	thisWindow, _ := ui.GetWindow(w.name)
	thisWindow.ClearWindow()

	lines := []string{}

	i := 0
	cursorPosition := 0
	curr := w.stateService.state.collection.head
	for curr != nil {
		lines = append(lines, curr.Name)

		if curr.Id == w.stateService.state.collection.selected.Id {
			cursorPosition = i
		}

		curr = curr.Next
		i += 1
	}

	thisWindow.WriteLines(lines)

	err := thisWindow.SetCursor(0, cursorPosition)
	if err != nil {
		log.Panic(err)
	}
}

// =============== ACTIONS ======================

func (rw *RequestsWindow) doRequest(ui *ui.UI) {
	if rw.stateService.state.collection.selected == nil {
		rw.stateService.state.alertMessage = "Create a request first"
		win, _ := ui.GetWindow("AlertWindow")
		win.OpenWindow()
		return
	}

	err := rw.stateService.ExecuteRequest()
	if err != nil {
		rw.stateService.state.alertMessage = err.Error()

		alertWindow, _ := ui.GetWindow("AlertWindow")
		alertWindow.OpenWindow()
	}

	win, _ := ui.GetWindow("ResponseWindow")
	win.Window.ReloadContent(ui, win)
}

func (rw *RequestsWindow) navigateDown(ui *ui.UI) {
	thisWindow, _ := ui.GetWindow(rw.name)

	ok := rw.stateService.SelectNext()
	if !ok {
		return
	}

	rw.ReloadContent(ui, thisWindow)
	win, _ := ui.GetWindow("ResponseWindow")
	win.Window.ReloadContent(ui, win)
}

func (rw *RequestsWindow) navigateUp(ui *ui.UI) {
	thisWindow, _ := ui.GetWindow(rw.name)

	ok := rw.stateService.SelectPrev()
	if !ok {
		return
	}

	rw.ReloadContent(ui, thisWindow)
	win, errr := ui.GetWindow("ResponseWindow")
	if errr != nil {
		log.Panic("Pudim")
	}
	win.Window.ReloadContent(ui, win)
}

func (rw *RequestsWindow) deleteRequest(ui *ui.UI) {
	thisWindow, _ := ui.GetWindow(rw.name)

	rw.stateService.DeleteSelectedRequest()
	rw.ReloadContent(ui, thisWindow)

	win, _ := ui.GetWindow("ResponseWindow")
	win.Window.ReloadContent(ui, win)
	win, _ = ui.GetWindow("RequestDetailsWindow")
	win.Window.ReloadContent(ui, win)
}
