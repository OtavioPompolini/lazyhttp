package app

import (
	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/types"
)

type StateService struct {
	persistance database.PersistanceAdapter
	state       *State
}

type State struct {
	requests        []*types.Request
	selectedRequest *types.Request
	debuggerMode    bool
	// app configs here
}

func NewStateService(db database.PersistanceAdapter) *StateService {
	return &StateService{
		persistance: db,
		state:       loadState(db),
	}
}

func loadState(db database.PersistanceAdapter) *State {
	reqs := db.RequestRepository.GetRequests()

	selectedRequest := &types.Request{}
	if len(reqs) >= 0 {
		selectedRequest = reqs[0]
	}

	return &State{
		requests:        reqs,
		selectedRequest: selectedRequest,
	}
}

// func (m *Memory) ListRequests() []types.Request {
// 	return m.localStorage.ListRequests()
// }

func (ss StateService) CreateRequest(reqName string) *types.Request {
	// ss.unselectRequests()
	saved := ss.persistance.RequestRepository.CreateRequest(reqName)
	// saved.Selected = true
	ss.state.requests = append(ss.state.requests, saved)
	ss.state.selectedRequest = saved
	return saved
}

func (ss *StateService) UpdateRequest(r *types.Request) {
	ss.persistance.RequestRepository.UpdateRequest(r)
	ss.state.selectedRequest = r
}

func (ss *StateService) SelectNext() {
	for i, v := range ss.state.requests[:len(ss.state.requests)-1] {
		if v.Id == ss.state.selectedRequest.Id {
			ss.state.selectedRequest = ss.state.requests[i+1]
			break
		}
	}
}

func (ss *StateService) SelectPrev() {
	for i, v := range ss.state.requests[1:] {
		if v.Id == ss.state.selectedRequest.Id {
			ss.state.selectedRequest = ss.state.requests[i]
			break
		}
	}
}

// func (ss StateService) unselectRequests() {
// 	for _, v := range ss.state.requests {
// 		v.Selected = false
// 	}
// }

// func (m *Memory) IsEmpty() bool {
// 	return m.localStorage.IsEmpty()
// }
//
// func (m *Memory) GetSelectedRequest() *types.Request {
// 	return m.localStorage.GetSelectedRequest()
// }
//
//
// func (m *Memory) CreateResponse(r *types.Request) {
// 	m.localStorage.UpdateSelectedRequest(r)
// }

// func (s *State)
