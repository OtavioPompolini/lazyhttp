package app

import (
	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/memory"
	"github.com/OtavioPompolini/project-postman/internal/model"
	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type RequestDetailsWindow struct {
	name       string
	x, y       int
	w, h       int
	body       string
	isActive   bool
	isSelected bool
	memory     *memory.Memory
}

func NewRequestDetailsWindow(GUI *ui.UI, memory *memory.Memory) *ui.Window {
	_, b := GUI.Size()
	return ui.NewWindow(
		&RequestDetailsWindow{
			name:       "RequestDetailsWindow",
			x:          50,
			y:          0,
			h:          b - 1,
			w:          80,
			isSelected: false,
			memory:     memory,
		},
		true,
	)
}

func (w RequestDetailsWindow) Name() string {
	return w.name
}

func (w *RequestDetailsWindow) Setup(ui ui.UI, v ui.Window) {
	v.SetTitle(v.Window.Name())
	v.SetEditable(true)
}

func (w *RequestDetailsWindow) Update(ui ui.UI, v ui.Window) {
	if !w.isSelected {
		v.ClearWindow()
		v.WriteLn(w.memory.GetSelectedRequest().Body)
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

func (w *RequestDetailsWindow) SetKeybindings(ui *ui.UI) error {

	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEsc, func(g *gocui.Gui, v *gocui.View) error {
		win, err := ui.GetWindow("RequestsWindow")
		if err != nil {
			return err
		}
		ui.SelectWindow(win)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *RequestDetailsWindow) OnDeselect(ui ui.UI, v ui.Window) error {
	selected := w.memory.GetSelectedRequest()
	w.memory.UpdateRequest(
		&model.Request{
			Id:   selected.Id,
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
