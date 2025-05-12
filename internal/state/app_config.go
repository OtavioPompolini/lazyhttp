package state

import (
	"strconv"

	"github.com/OtavioPompolini/project-postman/internal/database"
)

const (
	SHOW_RESPONSE_HEADERS string = "showResponseHeaders"
)

type AppConfig struct {
	showResponseHeaders bool
	debuggerMode        bool

	configRepo database.ConfigRepository
}

func NewAppConfig(db database.PersistanceAdapter) *AppConfig {
	configMap := db.ConfigRepository.GetConfig()
	srh, _ := strconv.ParseBool(configMap[SHOW_RESPONSE_HEADERS])

	return &AppConfig{
		showResponseHeaders: srh,
		debuggerMode:        false,

		configRepo: db.ConfigRepository,
	}
}

func (ac *AppConfig) ShowResponseHeaders() bool {
	return ac.showResponseHeaders
}

func (ac *AppConfig) ToggleShowResponseHeaders() {
	ac.showResponseHeaders = !ac.showResponseHeaders
	ac.configRepo.Update(SHOW_RESPONSE_HEADERS, strconv.FormatBool(ac.showResponseHeaders))
}
