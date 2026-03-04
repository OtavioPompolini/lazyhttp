package app

import (
	"time"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type AlertWindow struct {
	name           string
	pendingMessage string
	windowPosition ui.WindowPosition
	thisWindow     *ui.Window
	thisUI         *ui.UI
}

func NewAlertWindow(GUI *ui.UI, eb *state.EventBus) *ui.Window {
	aw := &AlertWindow{
		name:   "AlertWindow",
		thisUI: GUI,
		windowPosition: ui.NewWindowPosition(
			25, 49, 50, 10,
			ui.RELATIVE, ui.RELATIVE, ui.RELATIVE, ui.FIXED,
		),
	}

	windowRef := ui.NewWindow(aw, false)
	aw.thisWindow = windowRef

	eb.Subscribe(state.AlertMessage, aw.onAlert())
	return windowRef
}

func (aw *AlertWindow) onAlert() func(e state.Event) {
	return func(e state.Event) {
		message, ok := e.Data.(string)
		if !ok {
			return
		}
		aw.pendingMessage = message
		aw.thisWindow.OpenWindow()
	}
}

func (aw *AlertWindow) Setup(ui *ui.UI) {
	aw.thisWindow.Write(aw.pendingMessage)
	go func() {
		time.Sleep(2 * time.Second)
		aw.thisUI.Update(func() {
			aw.thisUI.DeleteWindowByName(aw.name)
		})
	}()
}

func (aw *AlertWindow) Update(ui ui.UI) {}

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

func (aw *AlertWindow) SetKeybindings(ui *ui.UI) error {
	return nil
}
