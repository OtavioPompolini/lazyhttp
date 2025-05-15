package app

// import (
// 	"strings"
//
// 	"github.com/awesome-gocui/gocui"
//
// 	"github.com/OtavioPompolini/project-postman/internal/state"
// 	"github.com/OtavioPompolini/project-postman/internal/ui"
// )
//
// type CreateRequestWindow struct {
// 	name           string
// 	newReqName     string
// 	stateService   *state.StateService
// 	windowPosition ui.WindowPosition
// }
//
// func NewCreateRequestWindow(GUI *ui.UI, stateService *state.StateService) *ui.Window {
// 	return ui.NewWindow(
// 		&CreateRequestWindow{
// 			name:         "CreateRequestWindow",
// 			stateService: stateService,
// 			windowPosition: ui.NewWindowPosition(
// 				25, 49, 50, 2,
// 				ui.RELATIVE,
// 				ui.RELATIVE,
// 				ui.RELATIVE,
// 				ui.FIXED,
// 			),
// 		},
// 		false,
// 	)
// }
//
// func (w CreateRequestWindow) Name() string {
// 	return w.name
// }
//
// func (w *CreateRequestWindow) Setup(ui *ui.UI, v *ui.Window) {
// 	ui.SelectWindow(v)
// 	v.SetHightlight(true)
// 	v.SetEditable(true)
// 	v.SetTitle("New request name:")
// }
//
// func (w *CreateRequestWindow) Update(ui ui.UI, v ui.Window) {
// 	w.newReqName = strings.TrimSpace(v.GetWindowContent())
// }
//
// func (w *CreateRequestWindow) Size() ui.WindowPosition {
// 	return w.windowPosition
// }
//
// func (w *CreateRequestWindow) SetKeybindings(ui *ui.UI, win *ui.Window) error {
// 	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
// 		return w.createRequest(ui)
// 	}); err != nil {
// 		return err
// 	}
//
// 	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEsc, func(g *gocui.Gui, v *gocui.View) error {
// 		return w.closeWindow(ui)
// 	}); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (w *CreateRequestWindow) OnDeselect(ui ui.UI, v ui.Window) error {
// 	return nil
// }
//
// func (w *CreateRequestWindow) OnSelect(ui ui.UI, v ui.Window) error {
// 	return nil
// }
//
// //  ======= ACTIONS =======
//
// func (w *CreateRequestWindow) closeWindow(ui *ui.UI) error {
// 	ui.DeleteWindowByName(w.name)
//
// 	win, err := ui.GetWindow("RequestsWindow")
// 	if err != nil {
// 		return err
// 	}
//
// 	ui.SelectWindow(win)
//
// 	return nil
// }
//
// func (w *CreateRequestWindow) createRequest(ui *ui.UI) error {
// 	w.stateService.RequestsStateService.CreateRequest(w.newReqName)
// 	win, _ := ui.GetWindow("RequestsWindow")
//
// 	ui.DeleteWindowByName(w.name)
// 	ui.SelectWindow(win)
//
// 	// I REALLY, REALLY, REALLY DONT LIKE THIS
// 	win.Window.ReloadContent(ui, win)
//
// 	return nil
// }
//
// func (w *CreateRequestWindow) ReloadContent(ui *ui.UI, v *ui.Window) {
// }
