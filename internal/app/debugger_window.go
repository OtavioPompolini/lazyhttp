package app

import (
	"io"
	"log"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type DebuggerWindow struct {
	name           string
	windowPosition ui.WindowPosition
	visible        bool

	thisWindow *ui.Window
}

func NewDebuggerWindow(state *state.State) *ui.Window {
	debuggerWindow := &DebuggerWindow{
		name: "DebuggerWindow",
		windowPosition: ui.NewWindowPosition(
			0, 80, 100, 20,
			ui.RELATIVE, ui.RELATIVE, ui.RELATIVE, ui.RELATIVE,
		),
	}

	windowRef := ui.NewWindow(debuggerWindow, true)
	debuggerWindow.thisWindow = windowRef

	return windowRef
}

func (aw *DebuggerWindow) Setup(ui *ui.UI) {
	aw.thisWindow.AutoScroll()
	ui.SetDefaultOutput("DebuggerWindow", func(out io.Writer) {
		newLogger := log.New(out, "INFO: ", log.LstdFlags|log.Lshortfile)

		log.SetOutput(newLogger.Writer())
		log.SetPrefix(newLogger.Prefix())
		log.SetFlags(newLogger.Flags())
	})
}

func (aw *DebuggerWindow) Update(ui ui.UI) {}
func (aw *DebuggerWindow) Size() ui.WindowPosition {
	return aw.windowPosition
}

func (aw *DebuggerWindow) Name() string {
	return aw.name
}

func (aw *DebuggerWindow) SetKeybindings(ui *ui.UI) error {
	return nil
}

func (aw *DebuggerWindow) OnSelect(ui ui.UI, w ui.Window) error {
	return nil
}

func (aw *DebuggerWindow) OnDeselect(ui ui.UI, w ui.Window) error {
	return nil
}
