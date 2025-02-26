package app

import "github.com/OtavioPompolini/project-postman/ui"

type CreateRequestWindow struct {
	name string
	x, y int
	w, h int
	isActive bool
	onCreateRequest func(reqName string)
}

func NewCreateRequestWindow(GUI *ui.UI) *ui.Window {
	a, b := GUI.Size()
	return ui.NewWindow(
		&CreateRequestWindow{
			name:         "CreateRequestWindow",
			x:            a / 2,
			y:            b / 2,
			w:            20,
			h:            2,
			isActive: false,
		})
}

func (w CreateRequestWindow) Name() string {
	return w.name
}

func (w *CreateRequestWindow) Setup(v ui.Window) {
	// v.SelBgColor = gocui.ColorYellow
	// v.SetVimEditor()
	v.SetHightlight(true)
	v.SetEditable(true)
}

func (w *CreateRequestWindow) Update(v ui.Window) {
}

func (w *CreateRequestWindow) Size() (x, y, width, height int) {
	return w.x, w.y, w.x + w.w, w.y + w.h
}

func (w *CreateRequestWindow) IsActive() bool {
	return w.isActive
}

func (w *CreateRequestWindow) SetKeybindings(ui *ui.UI) error {
	return nil
}

func (w *CreateRequestWindow) OnDeselect() error {
	return nil
}

func (w *CreateRequestWindow) OnSelect() error {
	return nil
}
