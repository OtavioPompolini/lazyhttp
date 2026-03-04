package state

import (
	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/types"
)

type RequestEvent struct {
	Requests []*types.Request
	Pos      int
}

type RequestSystem struct {
	requests             map[int64][]*types.Request
	selectedCollectionId *int64
	pos                  int

	eventBus          *EventBus
	requestRepository database.RequestRepository
}

func newRequestSystem(db database.PersistanceAdapter, eb *EventBus) *RequestSystem {
	rs := &RequestSystem{
		requests:          make(map[int64][]*types.Request),
		requestRepository: db.RequestRepository,
		eventBus:          eb,
	}
	rs.wireEvents()
	return rs
}

func (rs *RequestSystem) wireEvents() {
	rs.eventBus.Subscribe(CollectionSelected, func(e Event) {
		event := e.Data.(CollectionSelectedEvent)
		rs.selectedCollectionId = &event.Collection.Id
		rs.pos = 0
		rs.publishRequestChanged()
	})
}

func (rs *RequestSystem) publishRequestChanged() {
	var reqs []*types.Request
	if rs.selectedCollectionId != nil {
		reqs = rs.requests[*rs.selectedCollectionId]
	}
	rs.eventBus.Publish(Event{
		Type: RequestChanged,
		Data: RequestEvent{
			Requests: reqs,
			Pos:      rs.pos,
		},
	})
}

func (rs *RequestSystem) init() {
	rs.requests = loadRequests(rs.requestRepository)
}

func loadRequests(requestRpository database.RequestRepository) map[int64][]*types.Request {
	requestsMap := make(map[int64][]*types.Request)
	requestsList := requestRpository.GetRequests()

	for _, v := range requestsList {
		_, ok := requestsMap[v.CollectionId]
		if !ok {
			requestsMap[v.CollectionId] = []*types.Request{v}
		} else {
			requestsMap[v.CollectionId] = append(requestsMap[v.CollectionId], v)
		}
	}

	return requestsMap
}

func (rs *RequestSystem) Create(reqName string) {
	saved := rs.requestRepository.Create(reqName, *rs.selectedCollectionId)
	currRequests, ok := rs.requests[*rs.selectedCollectionId]

	if !ok {
		rs.requests[*rs.selectedCollectionId] = []*types.Request{saved}
	} else {
		rs.requests[*rs.selectedCollectionId] = append(rs.requests[*rs.selectedCollectionId], saved)
	}

	rs.pos = len(currRequests) - 1
	rs.publishRequestChanged()
}

func (rs *RequestSystem) ListNames() []string {
	if rs.selectedCollectionId == nil || len(rs.requests[*rs.selectedCollectionId]) <= 0 {
		return []string{}
	}

	requestsList := []string{}

	for _, v := range rs.requests[*rs.selectedCollectionId] {
		requestsList = append(requestsList, v.Name)
	}

	return requestsList
}

func (rs *RequestSystem) CurrentPos() int {
	return rs.pos
}

func (rs *RequestSystem) SelectNext() {
	if rs.selectedCollectionId == nil {
		return
	}
	currRequests, ok := rs.requests[*rs.selectedCollectionId]
	if !ok {
		return
	}

	rs.pos = min(len(currRequests)-1, rs.pos+1)
	rs.publishRequestChanged()
}

func (rs *RequestSystem) SelectPrev() {
	if rs.selectedCollectionId == nil {
		return
	}
	_, ok := rs.requests[*rs.selectedCollectionId]
	if !ok {
		return
	}

	rs.pos = max(0, rs.pos-1)
	rs.publishRequestChanged()
}

func (rs *RequestSystem) Update(r *types.Request) {
	// c.requests[c.currColl][c.currReq].Body = r.Body
	// c.requestRepository.UpdateRequest(r)
}
