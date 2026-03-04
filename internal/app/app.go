package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/tui"
)

type App struct {
	program *tea.Program
	state   *state.State
}

func NewApp() (*App, error) {
	db, err := database.NewPersistanceAdapter()
	if err != nil {
		return nil, err
	}

	eb := state.NewEventBus()
	st := state.NewState(db, eb)

	rootModel := tui.NewRootModel(st, eb)
	p := tea.NewProgram(rootModel, tea.WithAltScreen())

	tui.NewBridge(p, eb)

	return &App{program: p, state: st}, nil
}

func (app *App) Run() error {
	// Run state.Init asynchronously so the program is running before events publish.
	go app.state.Init()
	_, err := app.program.Run()
	return err
}
