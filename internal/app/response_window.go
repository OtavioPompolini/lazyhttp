package app

import (
	"log"

	"github.com/OtavioPompolini/project-postman/internal/ui"
	"github.com/OtavioPompolini/project-postman/internal/utils"
	"github.com/awesome-gocui/gocui"
)

type ResponseWindow struct {
	name         string
	x, y         int
	w, h         int
	isActive     bool
	stateService StateService
}

func NewResponseWindow(GUI *ui.UI, ss StateService) *ui.Window {
	return ui.NewWindow(
		&ResponseWindow{
			name:         "ResponseWindow",
			x:            60,
			y:            0,
			w:            40,
			h:            80,
			isActive:     true,
			stateService: ss,
		},
		true,
	)
}

func (w ResponseWindow) Name() string {
	return w.name
}

func (w *ResponseWindow) Setup(ui *ui.UI, v *ui.Window) {
	v.SetTitle("Response")
	v.SetSelectedBgColor(gocui.ColorRed)
	// v.SetHightlight(true)
	v.Wrap(true)
	// v.SetHightlight(true)
	w.ReloadContent(ui, v)
}

func (w *ResponseWindow) Update(ui ui.UI, v ui.Window) {
}

func (w *ResponseWindow) Size() (x, y, width, height int) {
	return w.x, w.y, w.x + w.w, w.y + w.h
}

func (w *ResponseWindow) IsActive() bool {
	return w.isActive
}

func (w *ResponseWindow) SetKeybindings(ui *ui.UI, v *ui.Window) error {
	return nil
}

func (w *ResponseWindow) OnDeselect(ui ui.UI, v ui.Window) error {
	return nil
}

func (w *ResponseWindow) OnSelect(ui ui.UI, v ui.Window) error {
	return nil
}

func (w *ResponseWindow) ReloadContent(ui *ui.UI, v *ui.Window) {
	win, _ := ui.GetWindow(w.name)
	win.ClearWindow()
	_, _, _, d := win.Window.Size()
	win.SetCursor(0, d/2)

	if w.stateService.state.collection.selected == nil || w.stateService.state.collection.selected.ResponseHistory == nil || len(w.stateService.state.collection.selected.ResponseHistory) <= 0 {
		return
	}

	win.Write(w.stateService.state.collection.selected.ResponseHistory[0].Info)
	err := win.WriteFunc(utils.StringBeautify(w.stateService.state.collection.selected.ResponseHistory[0].Body))
	if err != nil {
		log.Print("Error while highlighting response body", err)
		w.stateService.state.alertMessage = "Error while highlighting response body"
		alertWindow, _ := ui.GetWindow("AlertWindow")
		alertWindow.OpenWindow()
		return
	}
}
