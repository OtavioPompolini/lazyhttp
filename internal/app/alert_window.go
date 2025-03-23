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
	a, b := GUI.Size()
	return ui.NewWindow(
		&AlertWindow{
			name:         "AlertWindow",
			x:            (a / 2) - 25,
			y:            b / 2,
			w:            50,
			h:            2,
			stateService: stateService,
		},
		false,
	)
}

func (aw *AlertWindow) Setup(ui ui.UI, w ui.Window) {
	thisWindow, _ := ui.GetWindow(aw.name)

	thisWindow.Write(aw.stateService.state.alertMessage)

	go func() {
		time.Sleep(5 * time.Second)

		ui.DeleteWindowByName(aw.name)
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
