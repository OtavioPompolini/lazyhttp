package app

import (
	"log"

	"github.com/OtavioPompolini/project-postman/internal/ui"
	"github.com/OtavioPompolini/project-postman/internal/utils"
	"github.com/awesome-gocui/gocui"
	"golang.design/x/clipboard"
)

type selectType string

type cursorPosition struct {
	x int
	y int
}

type selectContent struct {
	sType         selectType
	startPosition cursorPosition
}

type ResponseWindow struct {
	name             string
	stateService     StateService
	windowPosition   ui.WindowPosition
	selectContent    *selectContent
	highlightedLines []int
}

func NewResponseWindow(GUI *ui.UI, ss StateService) *ui.Window {
	return ui.NewWindow(
		&ResponseWindow{
			name:         "ResponseWindow",
			stateService: ss,
			windowPosition: ui.NewWindowPosition(
				60, 0, 40, 80,
				ui.RELATIVE, ui.RELATIVE, ui.RELATIVE, ui.RELATIVE,
			),
		},
		true,
	)
}

func (w ResponseWindow) Name() string {
	return w.name
}

func (w *ResponseWindow) Size() ui.WindowPosition {
	return w.windowPosition
}

func (w *ResponseWindow) Setup(ui *ui.UI, v *ui.Window) {
	v.SetTitle("Response")
	v.SetSelectedBgColor(gocui.ColorRed)
	// v.SetHightlight(true)
	v.Wrap(true)
	// v.SetHightlight(true)
	w.ReloadContent(ui, v)
}

func (w *ResponseWindow) Update(ui ui.UI, v ui.Window) {
}

func (w *ResponseWindow) SetKeybindings(ui *ui.UI, v *ui.Window) error {
	if err := ui.NewKeyBinding(w.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
		w.navigateDown(ui)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'k', func(g *gocui.Gui, v *gocui.View) error {
		w.navigateUp(ui)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'V', func(g *gocui.Gui, v *gocui.View) error {
		win, _ := ui.GetWindow(w.name)
		win.SetSelectedBgColor(gocui.ColorRed)
		x, y := win.Cursor()
		win.HighlightLine(y, true)
		w.highlightedLines = append(w.highlightedLines, y)
		w.selectContent = &selectContent{
			sType: "LINE",
			startPosition: cursorPosition{
				x: x,
				y: y,
			},
		}

		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'y', func(g *gocui.Gui, v *gocui.View) error {
		w.copySelectedToClipboard(ui)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEsc, func(g *gocui.Gui, v *gocui.View) error {
		w.deselectContent(ui)

		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), gocui.KeyF1, func(g *gocui.Gui, v *gocui.View) error {
		ui.SelectWindowByName("RequestsWindow")

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *ResponseWindow) OnDeselect(ui ui.UI, v ui.Window) error {
	ui.CursorVisible(false)
	return nil
}

func (w *ResponseWindow) OnSelect(ui ui.UI, v ui.Window) error {
	ui.CursorVisible(true)
	return nil
}

func (w *ResponseWindow) ReloadContent(ui *ui.UI, v *ui.Window) {
	win, _ := ui.GetWindow(w.name)
	win.ClearWindow()

	if w.stateService.state.collection.selected == nil || w.stateService.state.collection.selected.ResponseHistory == nil || len(w.stateService.state.collection.selected.ResponseHistory) <= 0 {
		return
	}

	win.Write(w.stateService.state.collection.selected.ResponseHistory[0].Info)
	err := win.WriteFunc(utils.StringBeautify(w.stateService.state.collection.selected.ResponseHistory[0].Body))
	if err != nil {
		log.Print("Error while highlighting response body", err)
		w.stateService.state.alertMessage = "Error while highlighting response body"
		alertWindow, _ := ui.GetWindow("AlertWindow")
		alertWindow.OpenWindow()
		return
	}
}

// ACTIONS

func (rw *ResponseWindow) navigateDown(ui *ui.UI) {
	thisWindow, _ := ui.GetWindow(rw.name)

	if rw.selectContent != nil {
		_, y := thisWindow.Cursor()
		if y < rw.selectContent.startPosition.y {
			thisWindow.HighlightLine(y, false)
			rw.highlightedLines = append(rw.highlightedLines, y)
		}
	}

	thisWindow.CursorDown()

	if rw.selectContent != nil {
		_, y := thisWindow.Cursor()
		if y > rw.selectContent.startPosition.y {
			thisWindow.HighlightLine(y, true)
			rw.highlightedLines = append(rw.highlightedLines, y)
		}
	}
}

func (rw *ResponseWindow) navigateUp(ui *ui.UI) {
	thisWindow, _ := ui.GetWindow(rw.name)

	if rw.selectContent != nil {
		_, y := thisWindow.Cursor()
		if y > rw.selectContent.startPosition.y {
			thisWindow.HighlightLine(y, false)
			rw.highlightedLines = append(rw.highlightedLines, y)
		}

	}

	thisWindow.CursorUp()

	if rw.selectContent != nil {
		_, y := thisWindow.Cursor()
		if y < rw.selectContent.startPosition.y {
			thisWindow.HighlightLine(y, true)
			rw.highlightedLines = append(rw.highlightedLines, y)
		}
	}
}

func (rw *ResponseWindow) deselectContent(ui *ui.UI) {
	thisWindow, _ := ui.GetWindow(rw.name)

	for _, line := range rw.highlightedLines {
		thisWindow.HighlightLine(line, false)
	}

	rw.selectContent = nil
	rw.highlightedLines = []int{}
}

// Create utils to copy string
func (rw *ResponseWindow) copySelectedToClipboard(ui *ui.UI) {
	if rw.selectContent == nil {
		log.Println("No content selected")
		// Alert screen
		return
	}

	thisWindow, _ := ui.GetWindow(rw.name)
	err := clipboard.Init()
	if err != nil {
		log.Println("Failed to access clipboard")
		// Alert screen
		return
	}

	cx, cy := thisWindow.Cursor()

	startLine := min(rw.selectContent.startPosition.y, cy)
	endLine := max(rw.selectContent.startPosition.y, cy)
	copyContent := ""

	for i := startLine; i <= endLine; i++ {
		copyContent += thisWindow.Line(i) + "\n"
	}

	clipboard.Write(clipboard.FmtText, []byte(copyContent))

	thisWindow.SetCursor(cx, min(rw.selectContent.startPosition.y, cy))
	rw.deselectContent(ui)
}
