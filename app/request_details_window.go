package app

import "github.com/OtavioPompolini/project-postman/ui"

type RequestDetailsWindow struct {
	name string
	x, y int
	w, h int
	body string
	isActive bool
	toggleCursor func(b bool)
}

func NewRequestDetailsWindow(GUI *ui.UI) *ui.Window {
	return ui.NewWindow(
		&RequestDetailsWindow{
			name: "RequestDetailsWindow",
			x:    21,
			y:    0,
			h:    40,
			w:    40,
			isActive: true,
			toggleCursor: GUI.SetCursor,
		})
}

func (w RequestDetailsWindow) Name() string {
	return w.name
}

func (w *RequestDetailsWindow) Setup(v ui.Window) {
	// v.SelBgColor = gocui.ColorYellow
	// v.SetVimEditor()
	v.SetTitle(v.Window.Name())
	v.SetHightlight(true)
	v.SetEditable(true)
}

func (w *RequestDetailsWindow) Update(v ui.Window) {
}

func (w *RequestDetailsWindow) Size() (x, y, width, height int) {
	return w.x, w.y, w.x + w.w, w.y + w.h
}

func (w *RequestDetailsWindow) IsActive() bool {
	return w.isActive
}

func (w *RequestDetailsWindow) SetKeybindings(ui ui.UI) error {
	return nil
}

func (w *RequestDetailsWindow) OnDeselect() error {
	w.toggleCursor(false)
	return nil
}

func (w *RequestDetailsWindow) OnSelect() error {
	w.toggleCursor(true)
	return nil
}

//Bindings

// if err := g.SetKeybinding(w.Name, gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
// 	ChangeView(g, LastViewId)
// 	return nil
// }); err != nil {
// 	log.Panicln(err)
// }
