package memory

import (
	"errors"
	"sort"

	"github.com/OtavioPompolini/project-postman/model"
)

type LocalMemory struct {
	requests    map[int64]model.Request //TODO: support collections (filesystem)
	requestsArr []model.Request
	selectedReq int64
	selectedPos int
}

func newLocalMemory() *LocalMemory {
	return &LocalMemory{
		requests: map[int64]model.Request{},
		requestsArr: []model.Request{},
	}
}

func (m *LocalMemory) AddRequest(r *model.Request) {
	m.requests[r.Id] = *r

	m.reloadList()
}

func (m *LocalMemory) addRequests(r *map[int64]model.Request) {
	m.requests = *r
	m.reloadList()
}

func (m *LocalMemory) SelectNext() {
	if len(m.requestsArr) <= 0 {
		return
	}

	m.selectedPos = min(m.selectedPos+1, len(m.requestsArr)-1)
	m.selectedReq = m.requestsArr[m.selectedPos].Id
}

func (m *LocalMemory) SelectPrev() {
	if len(m.requestsArr) <= 0 {
		return
	}

	m.selectedPos = max(m.selectedPos-1, 0)
	m.selectedReq = m.requestsArr[m.selectedPos].Id
}

func (m *LocalMemory) ListRequests() []model.Request {
	return m.requestsArr
}

func (m *LocalMemory) IsEmpty() bool {
	if len(m.requestsArr) <= 0 {
		return true
	}

	return false
}

func (m *LocalMemory) GetSelectedRequest() model.Request {
	req, ok := m.requests[m.selectedReq]
	if ok {
		return req
	}

	return model.Request{}
}

func (m *LocalMemory) UpdateSelectedRequest(r *model.Request) {
	saved, _ := m.requests[r.Id]
	m.requests[r.Id] = model.Request{
		Id: r.Id,
		Name: saved.Name,
		Body: r.Body,
	}

	m.reloadList()
}

func (m *LocalMemory) reloadList() {
	requestsArr := []model.Request{}

	for _, req := range m.requests {
		requestsArr = append(requestsArr, req)
	}

	sort.Slice(requestsArr, func(i, j int) bool {
		return requestsArr[i].Id < requestsArr[j].Id
	})

	m.requestsArr = requestsArr
}

func (m *LocalMemory) SetSelected(selectedId int64) error {
	_, ok := m.requests[selectedId]
	if ok {
		return nil
	}

	return errors.New("Request not found")
}
