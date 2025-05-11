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
	AppConfig        AppConfig
	alertMessage     string

	// I dont like the name variables, think another name
	// variables    map[string]types.Variable
	// app configs here

}

func NewState(db database.PersistanceAdapter) *State {
	reqs := db.RequestRepository.GetRequests()
	loadResponses(db, reqs)

	return &State{
		CollectionSystem: newCollectionSystem(db),
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

func (ss *State) AlertMessage() string {
	return ss.alertMessage
}

func (ss *State) SetAlertMessage(am string) {
	ss.alertMessage = am
}
