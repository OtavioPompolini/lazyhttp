package app

import (
	"time"

	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type AlertWindow struct {
	name           string
	stateService   StateService
	windowPosition ui.WindowPosition
}

func NewAlertWindow(GUI *ui.UI, stateService StateService) *ui.Window {
	return ui.NewWindow(
		&AlertWindow{
			name:         "AlertWindow",
			stateService: stateService,
			windowPosition: ui.NewWindowPosition(
				25, 49, 50, 10,
				ui.RELATIVE, ui.RELATIVE, ui.RELATIVE, ui.FIXED,
			),
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

func (aw *AlertWindow) Size() ui.WindowPosition {
	return aw.windowPosition
}

func (aw *AlertWindow) Name() string {
	return aw.name
}

func (aw *AlertWindow) SetKeybindings(ui *ui.UI, w *ui.Window) error {
	return nil
}

func (aw *AlertWindow) ReloadContent(ui *ui.UI, w *ui.Window) {
}
