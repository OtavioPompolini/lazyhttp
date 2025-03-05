package app

import (
	"strings"

	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/memory"
	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type CreateRequestWindow struct {
	name       string
	x, y       int
	w, h       int
	isActive   bool
	newReqName string
	memory     *memory.Memory
}

func NewCreateRequestWindow(GUI *ui.UI, mem *memory.Memory) *ui.Window {
	a, b := GUI.Size()
	return ui.NewWindow(
		&CreateRequestWindow{
			name:     "CreateRequestWindow",
			x:        (a / 2) - 25,
			y:        b / 2,
			w:        50,
			h:        2,
			isActive: false,
			memory:   mem,
		},
		false,
	)
}

func (w CreateRequestWindow) Name() string {
	return w.name
}

func (w *CreateRequestWindow) Setup(ui ui.UI, v ui.Window) {
	ui.SelectWindow(&v)
	v.SetHightlight(true)
	v.SetEditable(true)
}

func (w *CreateRequestWindow) Update(ui ui.UI, v ui.Window) {
	w.newReqName = strings.TrimSpace(v.GetWindowContent())
}

func (w *CreateRequestWindow) Size() (x, y, width, height int) {
	return w.x, w.y, w.x + w.w, w.y + w.h
}

func (w *CreateRequestWindow) IsActive() bool {
	return w.isActive
}

func (w *CreateRequestWindow) SetKeybindings(ui *ui.UI) error {
	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
		saved := w.memory.CreateRequest(w.newReqName)

		win, err := ui.GetWindow("CreateRequestWindow")
		if err != nil {
			return err
		}

		err = ui.DeleteWindow(win)
		if err != nil {
			return err
		}

		win, err = ui.GetWindow("RequestsWindow")
		if err != nil {
			return err
		}

		win.WriteLn(saved.Name)
		ui.SelectWindow(win)
		win.SetCursor(1, 0)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *CreateRequestWindow) OnDeselect(ui ui.UI, v ui.Window) error {
	return nil
}

func (w *CreateRequestWindow) OnSelect(ui ui.UI, v ui.Window) error {
	return nil
}
