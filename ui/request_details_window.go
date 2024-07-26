package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type RequestDetailsWindow struct {
	name string
	x, y int
	w, h int
	body string
}

func NewRequestDetailsWindow(ui *UI) *View {
	return &View{
		Window: &RequestDetailsWindow{
			name: "RequestDetailsWindow",
			x:    21,
			y:    0,
			h:    40,
			w:    40,
		},
	}
}

func (w RequestDetailsWindow) Name() string {
	return w.name
}

func (w *RequestDetailsWindow) Setup(g *gocui.Gui, v *gocui.View) {
	// v.SelBgColor = gocui.ColorYellow
	v.Highlight = true
}

func (w *RequestDetailsWindow) Update(g *gocui.Gui, v *gocui.View) {
	if g.CurrentView().Name() == w.name {
		v.Editable = true
		w.body = v.Buffer()
	}

	v.Clear()
	fmt.Fprint(v, w.body)
}

func (w *RequestDetailsWindow) Size() (x, y, width, height int) {
	return w.x, w.y, w.x + w.w, w.y + w.h
}

func (w *RequestDetailsWindow) SetKeybindings(ui UI) error {
	return nil
}

//Bindings

// if err := g.SetKeybinding(w.Name, gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
// 	ChangeView(g, LastViewId)
// 	return nil
// }); err != nil {
// 	log.Panicln(err)
// }
