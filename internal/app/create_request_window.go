package app

import (
	"strings"

	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type CreateRequestWindow struct {
	name          string
	newReqName    string
	requestSystem *state.RequestSystem
	windowPosition ui.WindowPosition
	thisWindow     *ui.Window
}

func NewCreateRequestWindow(GUI *ui.UI, st *state.State) *ui.Window {
	crw := &CreateRequestWindow{
		name:          "CreateRequestWindow",
		requestSystem: st.RequestSystem,
		windowPosition: ui.NewWindowPosition(
			25, 49, 50, 2,
			ui.RELATIVE, ui.RELATIVE, ui.RELATIVE, ui.FIXED,
		),
	}

	windowRef := ui.NewWindow(crw, false)
	crw.thisWindow = windowRef
	return windowRef
}

func (w CreateRequestWindow) Name() string {
	return w.name
}

func (w *CreateRequestWindow) Setup(ui *ui.UI) {
	ui.SelectWindow(w.thisWindow)
	w.thisWindow.SetHightlight(true)
	w.thisWindow.SetEditable(true)
	w.thisWindow.SetTitle("New request name:")
}

func (w *CreateRequestWindow) Update(ui ui.UI) {
	w.newReqName = strings.TrimSpace(w.thisWindow.GetWindowContent())
}

func (w *CreateRequestWindow) Size() ui.WindowPosition {
	return w.windowPosition
}

func (w *CreateRequestWindow) SetKeybindings(ui *ui.UI) error {
	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
		return w.createRequest(ui)
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEsc, func(g *gocui.Gui, v *gocui.View) error {
		return w.closeWindow(ui)
	}); err != nil {
		return err
	}

	return nil
}

func (w *CreateRequestWindow) OnDeselect(ui ui.UI, v ui.Window) error {
	return nil
}

func (w *CreateRequestWindow) OnSelect(ui ui.UI, v ui.Window) error {
	return nil
}

func (w *CreateRequestWindow) closeWindow(ui *ui.UI) error {
	ui.DeleteWindowByName(w.name)

	win, err := ui.GetWindow("RequestsWindow")
	if err != nil {
		return err
	}

	ui.SelectWindow(win)
	return nil
}

func (w *CreateRequestWindow) createRequest(ui *ui.UI) error {
	w.requestSystem.Create(w.newReqName)

	win, _ := ui.GetWindow("RequestsWindow")
	ui.DeleteWindowByName(w.name)
	ui.SelectWindow(win)

	return nil
}
