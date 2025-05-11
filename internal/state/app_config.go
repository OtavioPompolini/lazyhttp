package state

import "strconv"

const (
	SHOW_RESPONSE_HEADERS string = "showResponseHeaders"
)

type AppConfig struct {
	showResponseHeaders bool
	debuggerMode        bool
}

func NewAppConfig(configMap map[string]string) *AppConfig {
	srh, _ := strconv.ParseBool(configMap[SHOW_RESPONSE_HEADERS])

	return &AppConfig{
		showResponseHeaders: srh,
		debuggerMode:        false,
	}
}

func (ac *AppConfig) ShowResponseHeaders() bool {
	return ac.showResponseHeaders
}
