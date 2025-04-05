package app

import (
	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/types"
	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type RequestDetailsWindow struct {
	name           string
	body           string
	isSelected     bool
	StateService   StateService
	windowPosition ui.WindowPosition
}

func NewRequestDetailsWindow(GUI *ui.UI, stateStateService StateService) *ui.Window {
	return ui.NewWindow(
		&RequestDetailsWindow{
			name:         "RequestDetailsWindow",
			isSelected:   false,
			StateService: stateStateService,
			windowPosition: ui.NewWindowPosition(
				20, 0, 40, 80,
				ui.RELATIVE, ui.RELATIVE, ui.RELATIVE, ui.RELATIVE,
			),
		},
		true,
	)
}

func (w RequestDetailsWindow) Name() string {
	return w.name
}

func (w *RequestDetailsWindow) Setup(ui *ui.UI, v *ui.Window) {
	v.SetTitle("Details")
	v.EnableKeybindingOnEdit(false)
	v.SetEditable(true)
	v.SetVimEditor()
}

func (w *RequestDetailsWindow) Update(ui ui.UI, v ui.Window) {
	if !w.isSelected {
		v.ClearWindow()
		if w.StateService.state.collection.selected != nil {
			v.Write(w.StateService.state.collection.selected.Body)
		}
	} else {
		w.body = v.GetWindowContent()
	}
}

func (w *RequestDetailsWindow) Size() ui.WindowPosition {
	return w.windowPosition
}

func (w *RequestDetailsWindow) SetKeybindings(ui *ui.UI, win *ui.Window) error {
	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEsc, func(g *gocui.Gui, v *gocui.View) error {
		_, err := ui.SelectWindowByName("RequestsWindow")
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *RequestDetailsWindow) OnDeselect(ui ui.UI, v ui.Window) error {
	w.StateService.UpdateRequest(
		&types.Request{
			Id:   w.StateService.state.collection.selected.Id,
			Body: w.body,
		},
	)

	w.isSelected = false
	ui.CursorVisible(false)
	return nil
}

func (w *RequestDetailsWindow) OnSelect(ui ui.UI, v ui.Window) error {
	w.isSelected = true
	ui.CursorVisible(true)
	return nil
}

func (w *RequestDetailsWindow) ReloadContent(ui *ui.UI, v *ui.Window) {
}
