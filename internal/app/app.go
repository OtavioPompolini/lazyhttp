package app

import (
	"log"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/adrg/xdg"

	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/tui"
	"github.com/OtavioPompolini/project-postman/internal/vim"
)

type App struct {
	program *tea.Program
	state   *state.State
}

func NewApp() (*App, error) {
	db, err := database.NewPersistenceAdapter()
	if err != nil {
		return nil, err
	}

	eb := state.NewEventBus()
	st := state.NewState(db, eb)

	engine := vim.NewEngine()
	configPath := filepath.Join(xdg.ConfigHome, "lazyhttp", "init.lua")
	if err := vim.LoadConfig(configPath, &engine); err != nil {
		log.Printf("warn: vim config: %v", err)
	}

	rootModel := tui.NewRootModel(st, eb, engine)
	p := tea.NewProgram(rootModel)

	tui.NewBridge(p, eb)

	return &App{program: p, state: st}, nil
}

func (app *App) Run() error {
	// Run state.Init asynchronously so the program is running before events publish.
	go app.state.Init()
	_, err := app.program.Run()
	return err
}
