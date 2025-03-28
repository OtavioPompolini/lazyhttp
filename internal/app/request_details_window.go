package app

import (
	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/types"
	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type RequestDetailsWindow struct {
	name         string
	x, y         int
	w, h         int
	body         string
	isActive     bool
	isSelected   bool
	StateService StateService
}

func NewRequestDetailsWindow(GUI *ui.UI, stateStateService StateService) *ui.Window {
	return ui.NewWindow(
		&RequestDetailsWindow{
			name:         "RequestDetailsWindow",
			x:            20,
			y:            0,
			h:            80,
			w:            40,
			isSelected:   false,
			StateService: stateStateService,
		},
		true,
	)
}

func (w RequestDetailsWindow) Name() string {
	return w.name
}

func (w *RequestDetailsWindow) Setup(ui *ui.UI, v *ui.Window) {
	v.SetTitle("Details")
	v.SetEditable(true)
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

func (w *RequestDetailsWindow) Size() (x, y, width, height int) {
	return w.x, w.y, w.x + w.w, w.y + w.h
}

func (w *RequestDetailsWindow) IsActive() bool {
	return w.isActive
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
	ui.SetCursor(false)
	return nil
}

func (w *RequestDetailsWindow) OnSelect(ui ui.UI, v ui.Window) error {
	w.isSelected = true
	ui.SetCursor(true)
	return nil
}

func (w *RequestDetailsWindow) ReloadContent(ui *ui.UI, v *ui.Window) {
}
