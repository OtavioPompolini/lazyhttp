package app

import (
	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type VariablesWindow struct {
	stateService   StateService
	name           string
	windowPosition ui.WindowPosition
}

func NewVariablesWindow(GUI *ui.UI, stateService StateService) *ui.Window {
	return ui.NewWindow(
		&VariablesWindow{
			stateService: stateService,
			name:         "VariablesWindow",
			windowPosition: ui.NewWindowPosition(
				0, 40, 20, 40,
				ui.RELATIVE, ui.RELATIVE, ui.RELATIVE, ui.RELATIVE,
			),
		},
		true,
	)
}

func (w *VariablesWindow) Name() string {
	return w.name
}

func (w *VariablesWindow) Setup(ui *ui.UI, v *ui.Window) {
	v.SetTitle("Variables")
	w.ReloadContent(ui, v)
}

func (w *VariablesWindow) Update(ui ui.UI, v ui.Window) {
}

func (w *VariablesWindow) Size() ui.WindowPosition {
	return w.windowPosition
}

func (w *VariablesWindow) SetKeybindings(ui *ui.UI, win *ui.Window) error {
	if err := ui.NewKeyBinding(w.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
		// w.navigateDown(ui)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'k', func(g *gocui.Gui, v *gocui.View) error {
		// w.navigateUp(ui)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'D', func(g *gocui.Gui, v *gocui.View) error {
		// w.deleteRequest(ui)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'n', func(g *gocui.Gui, v *gocui.View) error {
		// win, _ := ui.GetWindow("CreateRequestWindow")
		// win.OpenWindow()
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

// Doest make sense to pass *ui.Window in this method
func (w *VariablesWindow) ReloadContent(ui *ui.UI, v *ui.Window) {
	// thisWindow, _ := ui.GetWindow(w.name)
	// thisWindow.ClearWindow()

	// lines := []string{}

	// i := 0
	// cursorPosition := 0
	// curr := w.stateService.state.collection.head
	// for curr != nil {
	// 	lines = append(lines, curr.Name)

	// 	if curr.Id == w.stateService.state.collection.selected.Id {
	// 		cursorPosition = i
	// 	}

	// 	curr = curr.Next
	// 	i += 1
	// }

	// thisWindow.WriteLines(lines)

	// err := thisWindow.SetCursor(0, cursorPosition)
	// if err != nil {
	// 	log.Panic(err)
	// }
}
