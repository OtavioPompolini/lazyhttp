package app

import (
	"strings"

	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type CreateRequestWindow struct {
	name         string
	x, y         int
	w, h         int
	isActive     bool
	newReqName   string
	stateService StateService
}

func NewCreateRequestWindow(GUI *ui.UI, stateStateService StateService) *ui.Window {
	return ui.NewWindow(
		&CreateRequestWindow{
			name:         "CreateRequestWindow",
			x:            25,
			y:            49,
			w:            50,
			h:            2,
			isActive:     false,
			stateService: stateStateService,
		},
		false,
	)
}

func (w CreateRequestWindow) Name() string {
	return w.name
}

func (w *CreateRequestWindow) Setup(ui *ui.UI, v *ui.Window) {
	ui.SelectWindow(v)
	v.SetHightlight(true)
	v.SetEditable(true)
	v.SetTitle("Create Request:")
}

func (w *CreateRequestWindow) Update(ui ui.UI, v ui.Window) {
	w.newReqName = strings.TrimSpace(v.GetWindowContent())
}

func (w *CreateRequestWindow) Size() (x, y, width, height int) {
	return w.x, w.y, w.x + w.w, w.y + w.h
}

func (w *CreateRequestWindow) IsActive() bool {
	return w.isActive
}

func (w *CreateRequestWindow) SetKeybindings(ui *ui.UI, win *ui.Window) error {
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

//  ======= ACTIONS =======

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
	w.stateService.CreateRequest(w.newReqName)
	win, _ := ui.GetWindow("RequestsWindow")

	ui.DeleteWindowByName(w.name)
	ui.SelectWindow(win)

	// I REALLY, REALLY, REALLY DONT LIKE THIS
	win.Window.ReloadContent(ui, win)

	return nil
}

func (w *CreateRequestWindow) ReloadContent(ui *ui.UI, v *ui.Window) {
}
