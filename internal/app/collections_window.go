package app

import (
	"math/rand/v2"
	"strconv"

	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/ui"
)

// Why not to save a ref to ui.GUI and ui.Window on every Window implementation?????
// Since every F method needs both params...
type CollectionsWindow struct {
	name           string
	windowPosition ui.WindowPosition

	collectionSystem *state.CollectionSystem
	thisWindow       *ui.Window
}

func NewCollectionWindow(GUI *ui.UI, state *state.State) *ui.Window {
	collectionsWindow := &CollectionsWindow{
		collectionSystem: state.CollectionSystem,
		name:             "CollectionsWindow",
		windowPosition: ui.NewWindowPosition(
			0, 0, 20, 40,
			ui.RELATIVE, ui.RELATIVE, ui.RELATIVE, ui.RELATIVE,
		),
	}

	state.CollectionSystem.Subscribe(collectionsWindow)

	windowRef := ui.NewWindow(
		collectionsWindow,
		true,
	)

	collectionsWindow.thisWindow = windowRef
	return windowRef
}

func (w *CollectionsWindow) Name() string {
	return w.name
}

func (w *CollectionsWindow) Setup(ui *ui.UI, v *ui.Window) {
	ui.SelectWindow(v)
	v.SetTitle("Collections")
	v.SetSelectedBgColor(gocui.ColorRed)
	v.SetHightlight(true)
}

func (w *CollectionsWindow) Update(ui ui.UI, v ui.Window) {
}

func (w *CollectionsWindow) Size() ui.WindowPosition {
	return w.windowPosition
}

func (w *CollectionsWindow) SetKeybindings(ui *ui.UI, win *ui.Window) error {
	if err := ui.NewKeyBinding(w.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
		w.collectionSystem.SelectNext()
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'k', func(g *gocui.Gui, v *gocui.View) error {
		w.collectionSystem.SelectPrev()
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'J', func(g *gocui.Gui, v *gocui.View) error {
		w.collectionSystem.SwapPosition(state.DIRECTION_DOWN)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'K', func(g *gocui.Gui, v *gocui.View) error {
		w.collectionSystem.SwapPosition(state.DIRECTION_UP)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'n', func(g *gocui.Gui, v *gocui.View) error {
		w.collectionSystem.NewCollection("pudim" + strconv.FormatInt(int64(rand.IntN(100)), 10))
		// win, _ := ui.GetWindow("CreateRequestWindow")
		// win.OpenWindow()
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), '2', func(g *gocui.Gui, v *gocui.View) error {
		ui.SelectWindowByName("RequestsWindow")
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
		// w.selectCollection()
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *CollectionsWindow) OnDeselect(ui ui.UI, v ui.Window) error {
	return nil
}

func (w *CollectionsWindow) OnSelect(ui ui.UI, v ui.Window) error {
	return nil
}

// =============== OBSERVER ACTIONS ==================

func (cw *CollectionsWindow) OnUpdateCollection() {
	cw.thisWindow.ClearWindow()
	cw.thisWindow.WriteLines(cw.collectionSystem.ListNames())
	cw.thisWindow.SetCursor(0, cw.collectionSystem.CurrentPos())
}

// =============== ACTIONS ======================

func (cw *CollectionsWindow) navigateDown(ui *ui.UI) {
	cw.collectionSystem.SelectNext()

	// thisWindow, _ := ui.GetWindow(rw.name)
	//
	// ok := rw.stateService.SelectNext()
	// if !ok {
	// 	return
	// }
	//
	// rw.ReloadContent(ui, thisWindow)
	// win, _ := ui.GetWindow("ResponseWindow")
	// win.Window.ReloadContent(ui, win)
}

func (rw *CollectionsWindow) navigateUp(ui *ui.UI) {
	// thisWindow, _ := ui.GetWindow(rw.name)
	//
	// ok := rw.stateService.SelectPrev()
	// if !ok {
	// 	return
	// }
	//
	// rw.ReloadContent(ui, thisWindow)
	// win, _ := ui.GetWindow("ResponseWindow")
	// win.Window.ReloadContent(ui, win)
}

func (rw *CollectionsWindow) deleteRequest(ui *ui.UI) {
	// thisWindow, _ := ui.GetWindow(rw.name)
	//
	// rw.stateService.DeleteSelectedRequest()
	// rw.ReloadContent(ui, thisWindow)
	//
	// win, _ := ui.GetWindow("ResponseWindow")
	// win.Window.ReloadContent(ui, win)
	// win, _ = ui.GetWindow("RequestDetailsWindow")
	// win.Window.ReloadContent(ui, win)
}
