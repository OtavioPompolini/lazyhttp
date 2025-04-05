package app

import (
	"io"
	"log"

	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type DebuggerWindow struct {
	name           string
	stateService   StateService
	windowPosition ui.WindowPosition
}

func NewDebuggerWindow(stateService StateService) *ui.Window {
	return ui.NewWindow(
		&DebuggerWindow{
			name:         "DebuggerWindow",
			stateService: stateService,
			windowPosition: ui.NewWindowPosition(
				0, 80, 100, 20,
				ui.RELATIVE,
				ui.RELATIVE,
				ui.RELATIVE,
				ui.RELATIVE,
			),
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

func (aw *DebuggerWindow) Size() ui.WindowPosition {
	return aw.windowPosition
}

func (aw *DebuggerWindow) Name() string {
	return aw.name
}

func (aw *DebuggerWindow) SetKeybindings(ui *ui.UI, w *ui.Window) error {
	return nil
}

func (aw *DebuggerWindow) ReloadContent(ui *ui.UI, w *ui.Window) {
}
