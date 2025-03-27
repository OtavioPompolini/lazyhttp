package app

import (
	"io"
	"log"

	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type DebuggerWindow struct {
	name         string
	x, y         int
	w, h         int
	stateService StateService
}

func NewDebuggerWindow(stateService StateService) *ui.Window {
	return ui.NewWindow(
		&DebuggerWindow{
			name:         "DebuggerWindow",
			x:            0,
			y:            80,
			w:            100,
			h:            20,
			stateService: stateService,
		},
		true,
	)
}

func (aw *DebuggerWindow) Setup(ui *ui.UI, w *ui.Window) {
	w.AutoScroll()
	ui.SetDefaultOutput("DebuggerWindow", func(out io.Writer) {
		newLogger := log.New(out, "INFO: ", log.LstdFlags|log.Lshortfile)

		log.SetOutput(newLogger.Writer())
		log.SetPrefix(newLogger.Prefix())
		log.SetFlags(newLogger.Flags())
	})
}

func (aw *DebuggerWindow) Update(ui ui.UI, w ui.Window) {}
func (aw *DebuggerWindow) OnSelect(ui ui.UI, w ui.Window) error {
	return nil
}

func (aw *DebuggerWindow) OnDeselect(ui ui.UI, w ui.Window) error {
	return nil
}

func (aw *DebuggerWindow) Size() (x, y, w, h int) {
	return aw.x, aw.y, aw.x + aw.w, aw.y + aw.h
}

func (aw *DebuggerWindow) Name() string {
	return aw.name
}

func (aw *DebuggerWindow) SetKeybindings(ui *ui.UI, w *ui.Window) error {
	return nil
}

func (aw *DebuggerWindow) ReloadContent(ui *ui.UI, w *ui.Window) {
}
