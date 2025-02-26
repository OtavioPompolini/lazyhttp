package app

import (
	"github.com/OtavioPompolini/project-postman/memory"
	"github.com/OtavioPompolini/project-postman/request"
	"github.com/OtavioPompolini/project-postman/ui"
	"github.com/jroimartin/gocui"
)

type RequestDetailsWindow struct {
	name            string
	x, y            int
	w, h            int
	body            string
	isActive        bool
	isSelected      bool
	memory          *memory.Memory
	enableCursor    func(b bool)
	onUpdateRequest func(r *request.Request)
}

func NewRequestDetailsWindow(GUI *ui.UI, memory *memory.Memory, onUpdateRequest func(r *request.Request)) *ui.Window {
	return ui.NewWindow(
		&RequestDetailsWindow{
			name:            "RequestDetailsWindow",
			x:               21,
			y:               0,
			h:               40,
			w:               40,
			isActive:        true,
			isSelected:      false,
			memory:          memory,
			enableCursor:    GUI.SetCursor,
			onUpdateRequest: onUpdateRequest,
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

	// if err := ui.NewKeyBinding(w.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
	// 	win, err := ui.GetWindow(w.Name())
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	win.SetEditable(true)
	// 	win.SetTitle("EDITING")
	//
	// 	return nil
	// }); err != nil {
	// 	return err
	// }

	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEsc, func(g *gocui.Gui, v *gocui.View) error {
		win, err := ui.GetWindow(w.Name())
		if err != nil {
			return err
		}

		win.SetEditable(false)
		win.SetTitle(w.Name())

		return nil
	}); err != nil {
		return err
	}

	// if err := g.SetKeybinding(w.Name, gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
	// 	ChangeView(g, LastViewId)
	// 	return nil
	// }); err != nil {
	// 	log.Panicln(err)
	// }

	return nil
}

//Im thinking both this functions will need to pass Window as parameter. Im wanting to set Cursor on position 0
func (w *RequestDetailsWindow) OnDeselect() error {
	// onSaveBodyContent(w.body)
	selected := w.memory.GetSelectedRequest()
	w.onUpdateRequest(&request.Request{
		Id: selected.Id,
		Body: w.body,
	})
	w.isSelected = false
	w.enableCursor(false)
	return nil
}

func (w *RequestDetailsWindow) OnSelect() error {
	w.isSelected = true
	w.enableCursor(true)
	return nil
}

//Bindings
