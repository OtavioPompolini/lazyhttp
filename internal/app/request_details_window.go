package app

import (
	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/types"
	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type RequestDetailsWindow struct {
	name           string
	body           string
	isSelected     bool
	requestSystem  *state.RequestSystem
	currentRequest *types.Request
	windowPosition ui.WindowPosition
	thisWindow     *ui.Window
}

func NewRequestDetailsWindow(GUI *ui.UI, st *state.State, eb *state.EventBus) *ui.Window {
	rdw := &RequestDetailsWindow{
		name:          "RequestDetailsWindow",
		requestSystem: st.RequestSystem,
		windowPosition: ui.NewWindowPosition(
			20, 0, 40, 80,
			ui.RELATIVE, ui.RELATIVE, ui.RELATIVE, ui.RELATIVE,
		),
	}

	windowRef := ui.NewWindow(rdw, true)
	rdw.thisWindow = windowRef

	eb.Subscribe(state.RequestChanged, rdw.onRequestChanged())
	return windowRef
}

func (w *RequestDetailsWindow) onRequestChanged() func(e state.Event) {
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

		if !w.isSelected {
			w.thisWindow.ClearWindow()
			if w.currentRequest != nil {
				w.thisWindow.Write(w.currentRequest.Body)
			}
		}
	}
}

func (w RequestDetailsWindow) Name() string {
	return w.name
}

func (w *RequestDetailsWindow) Setup(ui *ui.UI) {
	w.thisWindow.SetTitle("Details")
	w.thisWindow.EnableKeybindingOnEdit(false)
	w.thisWindow.SetEditable(true)
}

func (w *RequestDetailsWindow) Update(ui ui.UI) {
	if w.isSelected {
		w.body = w.thisWindow.GetWindowContent()
	}
}

func (w *RequestDetailsWindow) Size() ui.WindowPosition {
	return w.windowPosition
}

func (w *RequestDetailsWindow) SetKeybindings(ui *ui.UI) error {
	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEsc, func(g *gocui.Gui, v *gocui.View) error {
		_, err := ui.SelectWindowByName("RequestsWindow")
		return err
	}); err != nil {
		return err
	}

	return nil
}

func (w *RequestDetailsWindow) OnDeselect(ui ui.UI, v ui.Window) error {
	if w.currentRequest != nil {
		w.requestSystem.Update(&types.Request{
			Id:   w.currentRequest.Id,
			Body: w.body,
		})
	}

	w.isSelected = false
	ui.CursorVisible(false)
	return nil
}

func (w *RequestDetailsWindow) OnSelect(ui ui.UI, v ui.Window) error {
	w.isSelected = true
	ui.CursorVisible(true)
	return nil
}
