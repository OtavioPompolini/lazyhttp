package app

import (
	"time"

	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type AlertWindow struct {
	name         string
	x, y         int
	w, h         int
	stateService StateService
}

func NewAlertWindow(GUI *ui.UI, stateService StateService) *ui.Window {
	return ui.NewWindow(
		&AlertWindow{
			name:         "AlertWindow",
			x:            25,
			y:            49,
			w:            50,
			h:            10,
			stateService: stateService,
		},
		false,
	)
}

func (aw *AlertWindow) Setup(ui *ui.UI, w *ui.Window) {
	thisWindow, _ := ui.GetWindow(aw.name)

	thisWindow.Write(aw.stateService.state.alertMessage)

	go func() {
		time.Sleep(2 * time.Second)
		ui.Update(
			func() {
				ui.DeleteWindowByName(aw.name)
			})
	}()
}

func (aw *AlertWindow) Update(ui ui.UI, w ui.Window) {}
func (aw *AlertWindow) OnSelect(ui ui.UI, w ui.Window) error {
	return nil
}

func (aw *AlertWindow) OnDeselect(ui ui.UI, w ui.Window) error {
	return nil
}

func (aw *AlertWindow) Size() (x, y, w, h int) {
	return aw.x, aw.y, aw.x + aw.w, aw.y + aw.h
}

func (aw *AlertWindow) Name() string {
	return aw.name
}

func (aw *AlertWindow) SetKeybindings(ui *ui.UI, w *ui.Window) error {
	return nil
}

func (aw *AlertWindow) ReloadContent(ui *ui.UI, w *ui.Window) {
}
