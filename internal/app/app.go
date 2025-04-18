package app

import (
	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/ui"
)

type App struct {
	GUI          *ui.UI
	debuggerMode bool
	state        State
}

func NewApp() (*App, error) {
	userInterface, err := ui.NewUI()
	if err != nil {
		return nil, err
	}

	db, err := database.NewPersistanceAdapter()
	if err != nil {
		return nil, err
	}

	stateService := NewStateService(db)

	app := &App{
		// persistanceAdapter: db,
		GUI:   userInterface,
		state: *stateService.state,
	}

	app.GUI.StartUI()
	app.GUI.AddWindow(NewDebuggerWindow(*stateService))
	app.GUI.AddWindow(NewResponseWindow(userInterface, *stateService))
	app.GUI.AddWindow(NewAlertWindow(userInterface, *stateService))
	app.GUI.AddWindow(NewRequestDetailsWindow(userInterface, *stateService))
	app.GUI.AddWindow(NewCreateRequestWindow(userInterface, *stateService))
	app.GUI.AddWindow(NewRequestsWindow(userInterface, *stateService))
	app.GUI.AddWindow(NewVariablesWindow(userInterface, *stateService))

	app.GUI.SetHightlight(true)
	app.GUI.SetFgColor(gocui.ColorGreen)
	app.GUI.SetSelectedFgColor(gocui.ColorYellow)

	if err := app.GUI.SetGlobalKeybindings(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Run() error {
	defer app.GUI.Close()
	if err := app.GUI.Start(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}
