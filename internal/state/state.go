package state

import (
	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/types"
)

// Observer pattern to update windows????
// type StateService struct {
// 	RequestsStateService RequestStateService
// 	ConfigStateService   ConfigStateService
//
// 	state *State
//
// 	// VariablesStateService VariablesStateService
// }

// func NewStateService(db database.PersistanceAdapter) *StateService {
// 	state := loadState(db)
//
// 	return &StateService{
// 		RequestsStateService: newRequestStateService(state, db),
// 		ConfigStateService:   newConfigStateService(state, db),
// 		state:                state,
// 	}
// }

type State struct {
	CollectionSystem *CollectionSystem
	RequestSystem    *RequestSystem

	// AppConfig          *AppConfig
	// NotificationSystem *NotificationSystem
	// WindowsStateManager *WindowsStateManager

	// I dont like the name variables, think another name
	// variables    map[string]types.Variable
}

type NotificationSystem struct {
	message string

	alertObservers []AlertNotificationObserver
}

type AlertNotificationObserver interface {
	SendAlertNotification(message string)
}

func NewState(db database.PersistanceAdapter) *State {
	// reqs := db.RequestRepository.GetRequests()
	// loadResponses(db, reqs)
	collectionSystem := newCollectionSystem(db)

	return &State{
		CollectionSystem: collectionSystem,
		RequestSystem:    newRequestSystem(db, &collectionSystem.selId),
		// NotificationSystem: newNotificationSystem(),
		// AppConfig:        NewAppConfig(db),
		// variables:  map[string]types.Variable{},
	}
}

func loadResponses(db database.PersistanceAdapter, reqs []*types.Request) {
	responses := db.ResponseRepository.GetAll()
	for _, v := range reqs {
		r, ok := responses[v.Id]
		if ok {
			v.ResponseHistory = r
		}
	}
}

func newNotificationSystem() *NotificationSystem {
	return &NotificationSystem{}
}

func (ns *NotificationSystem) createAlert(mes string) {
	for _, v := range ns.alertObservers {
		v.SendAlertNotification(mes)
	}
}

func (ns *NotificationSystem) subscribeAlert(obs AlertNotificationObserver) {
	ns.alertObservers = append(ns.alertObservers, obs)
}
