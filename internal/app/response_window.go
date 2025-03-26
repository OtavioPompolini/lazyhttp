package app

import (
	"github.com/OtavioPompolini/project-postman/internal/ui"
	"github.com/awesome-gocui/gocui"
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

func (w *ResponseWindow) Setup(ui *ui.UI, v *ui.Window) {
	v.SetTitle("Response:")
	v.SetSelectedBgColor(gocui.ColorRed)
	v.SetHightlight(true)
	v.Wrap(true)
	// v.SetHightlight(true)
	w.ReloadContent(ui, v)
}

func (w *ResponseWindow) Update(ui ui.UI, v ui.Window) {
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
	win, _ := ui.GetWindow(w.name)
	win.ClearWindow()
	_, _, _, d := win.Window.Size()
	win.SetCursor(0, d/2)

	if w.stateService.state.collection.selected != nil && w.stateService.state.collection.selected.ResponseHistory != nil && len(w.stateService.state.collection.selected.ResponseHistory) > 0 {
		win.Write(w.stateService.state.collection.selected.ResponseHistory[0].Info)
		win.WriteHighlight(w.stateService.state.collection.selected.ResponseHistory[0].Body)
	}
}
