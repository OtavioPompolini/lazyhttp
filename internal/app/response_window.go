package app

import (
	"log"

	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/types"
	"github.com/OtavioPompolini/project-postman/internal/ui"
	"github.com/OtavioPompolini/project-postman/internal/utils"
)

type ResponseWindow struct {
	name           string
	currentRequest *types.Request
	windowPosition ui.WindowPosition
	thisWindow     *ui.Window
}

func NewResponseWindow(GUI *ui.UI, st *state.State, eb *state.EventBus) *ui.Window {
	rw := &ResponseWindow{
		name: "ResponseWindow",
		windowPosition: ui.NewWindowPosition(
			60, 0, 40, 80,
			ui.RELATIVE, ui.RELATIVE, ui.RELATIVE, ui.RELATIVE,
		),
	}

	windowRef := ui.NewWindow(rw, true)
	rw.thisWindow = windowRef

	eb.Subscribe(state.RequestChanged, rw.onRequestChanged())
	return windowRef
}

func (w *ResponseWindow) onRequestChanged() func(e state.Event) {
	return func(e state.Event) {
		event, ok := e.Data.(state.RequestEvent)
		if !ok {
			return
		}

		if len(event.Requests) > 0 && event.Pos >= 0 && event.Pos < len(event.Requests) {
			w.currentRequest = event.Requests[event.Pos]
		} else {
			w.currentRequest = nil
		}

		w.reloadContent()
	}
}

func (w *ResponseWindow) reloadContent() {
	w.thisWindow.ClearWindow()

	if w.currentRequest == nil || len(w.currentRequest.ResponseHistory) == 0 {
		return
	}

	w.thisWindow.Write(w.currentRequest.ResponseHistory[0].Info)
	err := w.thisWindow.WriteFunc(utils.StringBeautify(w.currentRequest.ResponseHistory[0].Body))
	if err != nil {
		log.Print("Error while highlighting response body", err)
	}
}

func (w ResponseWindow) Name() string {
	return w.name
}

func (w *ResponseWindow) Size() ui.WindowPosition {
	return w.windowPosition
}

func (w *ResponseWindow) Setup(ui *ui.UI) {
	w.thisWindow.SetTitle("Response")
	w.thisWindow.SetSelectedBgColor(gocui.ColorRed)
	w.thisWindow.Wrap(true)
	w.reloadContent()
}

func (w *ResponseWindow) Update(ui ui.UI) {}

func (w *ResponseWindow) SetKeybindings(ui *ui.UI) error {
	if err := ui.NewKeyBinding(w.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
		w.thisWindow.CursorDown()
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'k', func(g *gocui.Gui, v *gocui.View) error {
		w.thisWindow.CursorUp()
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), gocui.KeyF1, func(g *gocui.Gui, v *gocui.View) error {
		ui.SelectWindowByName("RequestsWindow")
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *ResponseWindow) OnDeselect(ui ui.UI, v ui.Window) error {
	ui.CursorVisible(false)
	return nil
}

func (w *ResponseWindow) OnSelect(ui ui.UI, v ui.Window) error {
	ui.CursorVisible(true)
	return nil
}
