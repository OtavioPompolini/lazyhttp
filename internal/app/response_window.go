package app

import (
	"log"

	"github.com/OtavioPompolini/project-postman/internal/ui"
	"github.com/OtavioPompolini/project-postman/internal/utils"
	"github.com/awesome-gocui/gocui"
)

type ResponseWindow struct {
	name           string
	stateService   StateService
	windowPosition ui.WindowPosition
}

func NewResponseWindow(GUI *ui.UI, ss StateService) *ui.Window {
	return ui.NewWindow(
		&ResponseWindow{
			name:         "ResponseWindow",
			stateService: ss,
			windowPosition: ui.NewWindowPosition(
				60, 0, 40, 80,
				ui.RELATIVE, ui.RELATIVE, ui.RELATIVE, ui.RELATIVE,
			),
		},
		true,
	)
}

func (w ResponseWindow) Name() string {
	return w.name
}

func (w *ResponseWindow) Size() ui.WindowPosition {
	return w.windowPosition
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
