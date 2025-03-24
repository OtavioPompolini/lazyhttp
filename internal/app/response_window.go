package app

import (
	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type ResponseWindow struct {
	name         string
	x, y         int
	w, h         int
	isActive     bool
	newReqName   string
	stateService StateService
}

func NewResponseWindow(GUI *ui.UI, ss StateService) *ui.Window {
	a, b := GUI.Size()
	return ui.NewWindow(
		&ResponseWindow{
			name:         "ResponseWindow",
			x:            (a * 60 / 100) + 2,
			y:            0,
			w:            a*40/100 - 2,
			h:            b - 1,
			isActive:     true,
			stateService: ss,
		},
		true,
	)
}

func (w ResponseWindow) Name() string {
	return w.name
}

func (w *ResponseWindow) Setup(ui ui.UI, v ui.Window) {
	v.SetEditable(true)
	v.SetTitle("Response:")
	v.Wrap(true)
}

func (w *ResponseWindow) Update(ui ui.UI, v ui.Window) {
	v.ClearWindow()
	if w.stateService.state.collection.selected != nil {
		v.Write(w.stateService.state.collection.selected.LastResponse)
	}
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
}
