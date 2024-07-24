package views


// type RequestDetailView = ui.View
//
// func NewRequestDetailView() *RequestDetailView {
// 	return ui.NewView("request_detail")
// }

// func NewRequestDetailView(g *gocui.Gui) *ui.View {
// 	return ui.View.NewV
// 	// return &RequestDetailView{
// 	// 	Name:    RequestDetailViewName,
// 	// 	X:       20,
// 	// 	Y:       0,
// 	// 	Width:   40,
// 	// 	Height:  40,
// 	// 	Request: nil,
// 	// }
// }
// type RequestDetailView struct {
// 	Title         string
// 	Name          string
// 	X, Y          int
// 	Width, Height int
// 	Request       *types.Request
// }

// func (w *RequestDetailView) Layout(g *gocui.Gui) error {
//
// 	v, err := g.SetView(w.Name, w.X, w.Y, w.X+w.Width, w.Y+w.Height)
// 	if err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
//
// 		if err := g.SetKeybinding(w.Name, gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
// 			ChangeView(g, LastViewId)
// 			return nil
// 		}); err != nil {
// 			log.Panicln(err)
// 		}
// 	}
//
// 	if w.Request == nil {
// 		return nil
// 	}
//
// 	if g.CurrentView().Name() == w.Name {
// 		v.Editable = true
// 		w.Request.Body = v.Buffer()
// 		return nil
// 	}
//
// 	v.Clear()
// 	fmt.Fprint(v, w.Request.Body)
// 	return nil
// }

// func ChangeView(g *gocui.Gui, viewId string) error {
// 	LastViewId = g.CurrentView().Name()
// 	_, err := g.SetCurrentView(viewId)
//
// 	return err
// }
