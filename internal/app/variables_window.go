package app

import (
	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type VariablesWindow struct {
	name           string
	windowPosition ui.WindowPosition
	thisWindow     *ui.Window
}

func NewVariablesWindow(GUI *ui.UI, st *state.State) *ui.Window {
	vw := &VariablesWindow{
		name: "VariablesWindow",
		windowPosition: ui.NewWindowPosition(
			0, 40, 20, 40,
			ui.RELATIVE, ui.RELATIVE, ui.RELATIVE, ui.RELATIVE,
		),
	}

	windowRef := ui.NewWindow(vw, true)
	vw.thisWindow = windowRef
	return windowRef
}

func (w *VariablesWindow) Name() string {
	return w.name
}

func (w *VariablesWindow) Setup(ui *ui.UI) {
	w.thisWindow.SetTitle("Variables")
}

func (w *VariablesWindow) Update(ui ui.UI) {}

func (w *VariablesWindow) Size() ui.WindowPosition {
	return w.windowPosition
}

func (w *VariablesWindow) SetKeybindings(ui *ui.UI) error {
	if err := ui.NewKeyBinding(w.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'k', func(g *gocui.Gui, v *gocui.View) error {
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'D', func(g *gocui.Gui, v *gocui.View) error {
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'n', func(g *gocui.Gui, v *gocui.View) error {
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *VariablesWindow) OnDeselect(ui ui.UI, v ui.Window) error {
	v.SetHightlight(true)
	return nil
}

func (w *VariablesWindow) OnSelect(ui ui.UI, v ui.Window) error {
	v.SetHightlight(false)
	return nil
}
