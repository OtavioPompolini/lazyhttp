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
	collection   *Collection
	debuggerMode bool
	alertMessage string
	// app configs here
}

type Collection struct {
	head     *types.Request
	tail     *types.Request
	selected *types.Request
}

func NewStateService(db database.PersistanceAdapter) *StateService {
	return &StateService{
		persistance: db,
		state:       loadState(db),
	}
}

func NewCollection(requests []*types.Request) *Collection {
	var head *types.Request
	var tail *types.Request
	var prev *types.Request

	for i, req := range requests {
		if i == 0 {
			head = req
		}

		if i == len(requests)-1 {
			tail = req
		}

		req.Prev = prev
		if prev != nil {
			prev.Next = req
		}

		prev = req
	}

	return &Collection{
		head:     head,
		tail:     tail,
		selected: head,
	}
}

func loadState(db database.PersistanceAdapter) *State {
	reqs := db.RequestRepository.GetRequests()

	return &State{
		collection: NewCollection(reqs),
	}
}

// func (m *Memory) ListRequests() []types.Request {
// 	return m.localStorage.ListRequests()
// }

func (ss StateService) DeleteSelectedRequest() {
	selected := ss.state.collection.selected
	prev := ss.state.collection.selected.Prev
	next := ss.state.collection.selected.Next

	if prev != nil {
		prev.Next = selected.Next
		ss.state.collection.selected = prev
	}

	if next != nil {
		next.Prev = selected.Prev
		ss.state.collection.selected = next
	}

	ss.persistance.RequestRepository.DeleteRequest(selected.Id)
}

func (ss StateService) CreateRequest(reqName string) *types.Request {
	// ss.unselectRequests()
	saved := ss.persistance.RequestRepository.CreateRequest(reqName)
	// saved.Selected = true
	ss.state.collection.tail.Next = saved
	saved.Prev = ss.state.collection.tail
	ss.state.collection.selected = saved

	return saved
}

func (ss *StateService) UpdateRequest(r *types.Request) {
	ss.persistance.RequestRepository.UpdateRequest(r)
	ss.state.collection.selected.Body = r.Body
}

func (ss *StateService) SelectNext() {
	next := ss.state.collection.selected.Next
	if next != nil {
		ss.state.collection.selected = next
	}
}

func (ss *StateService) SelectPrev() {
	prev := ss.state.collection.selected.Prev
	if prev != nil {
		ss.state.collection.selected = prev
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
