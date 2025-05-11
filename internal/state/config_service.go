package state

import (
	"strconv"

	"github.com/OtavioPompolini/project-postman/internal/database"
)

type ConfigStateService struct {
	persistance database.PersistanceAdapter
	state       *State
}

func newConfigStateService(state *State, db database.PersistanceAdapter) ConfigStateService {
	return ConfigStateService{
		persistance: db,
		state:       state,
	}
}

func (css ConfigStateService) ToggleShowResponseHeader() {
	css.state.appConfig.showResponseHeaders = !css.state.appConfig.showResponseHeaders
	css.persistance.ConfigRepository.Save(SHOW_RESPONSE_HEADERS, strconv.FormatBool(css.state.appConfig.showResponseHeaders))
}
