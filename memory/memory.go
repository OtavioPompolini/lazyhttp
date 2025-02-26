package memory

import (
	"errors"
	"sort"

	"github.com/OtavioPompolini/project-postman/request"
)

type Memory struct {
	requests    map[int64]request.Request //TODO: support collections (filesystem)
	requestsArr []request.Request
	selectedReq int64
	selectedPos int
}

//Why adapter here??? Is it worth to pass adapter or only pass all memory elements?
//Thinking better, yes. adapter should be here. Its better to update memory and database on same place
//Prevent to unsync
func NewMemory(adapter request.Adapter) *Memory {
	requests := adapter.GetRequests()
	requestsArr := []request.Request{}

	for _, r := range *requests {
		requestsArr = append(requestsArr, r)
	}

	sort.Slice(requestsArr, func(i, j int) bool {
		return requestsArr[i].Id < requestsArr[j].Id
	})

	return &Memory{
		requests: *requests,
		requestsArr: requestsArr,
		selectedPos: 0,
	}
}

func (m *Memory) AddRequest(r *request.Request) {
	m.requests[r.Id] = *r

	m.reloadList()
}

func (m *Memory) SelectNext() {
	if len(m.requestsArr) == 0 {
		return
	}

	m.selectedPos += 1

	if m.selectedPos >= len(m.requestsArr) {
		m.selectedPos = len(m.requestsArr) - 1
	}

	m.selectedReq = m.requestsArr[m.selectedPos].Id
}

func (m *Memory) IsEmpty() bool {
	if len(m.requestsArr) <= 0 {
		return true
	}

	return false
}

func (m *Memory) SelectPrev() {
	if len(m.requestsArr) == 0 {
		return
	}

	m.selectedPos -= 1

	if m.selectedPos < 0 {
		m.selectedPos = 0
	}

	m.selectedReq = m.requestsArr[m.selectedPos].Id
}

func (m *Memory) ListRequests() []request.Request {
	return m.requestsArr
}

func (m *Memory) GetSelectedRequest() request.Request {
	req, ok := m.requests[m.selectedReq]
	if ok {
		return req
	}

	return request.Request{}
}

func (m *Memory) UpdateSelectedRequest(r *request.Request) {
	saved, _ := m.requests[r.Id]
	m.requests[r.Id] = request.Request{
		Id: r.Id,
		Name: saved.Name,
		Body: r.Body,
	}

	m.reloadList()
}

func (m *Memory) reloadList() {
	requestsArr := []request.Request{}

	for _, req := range m.requests {
		requestsArr = append(requestsArr, req)
	}

	sort.Slice(requestsArr, func(i, j int) bool {
		return requestsArr[i].Id < requestsArr[j].Id
	})

	m.requestsArr = requestsArr
}

// func (m *Memory) GetSelectedId() int64 {
// 	return m.selectedReq
// }

func (m *Memory) SetSelected(selectedId int64) error {
	_, ok := m.requests[selectedId]
	if ok {
		return nil
	}

	return errors.New("Request not found")
}
