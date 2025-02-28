package memory

import "github.com/OtavioPompolini/project-postman/internal/model"

// For now its kinda Oukay. I really don't know what will happen if the Memory
// grows in objects and this struct maybe will have a lot of methods implemented
// ITs starting to grow. Now we have responses to save too
type Memory struct {
	persistance  PersistanceAdapter
	localStorage *LocalMemory
}

// TODO: FUUUUUUUTURE accept more than only sqlite for persistance.
// Maybe a mysql db for shared collections?
func InitMemory() (*Memory, error) {
	db, err := initDatabase()
	if err != nil {
		return nil, err
	}

	localStorage := newLocalMemory()

	mem := &Memory{
		persistance:  db,
		localStorage: localStorage,
	}

	mem.loadLocalMemory()
	return mem, nil
}

func (m *Memory) loadLocalMemory() {
	reqs := m.persistance.GetRequests()
	m.localStorage.addRequests(reqs)
}

func (m *Memory) ListRequests() []model.Request {
	return m.localStorage.ListRequests()
}

func (m *Memory) CreateRequest(reqName string) *model.Request {
	saved := m.persistance.CreateRequest(reqName)
	m.localStorage.AddRequest(saved)
	return saved
}

func (m *Memory) SelectNext() {
	m.localStorage.SelectNext()
}

func (m *Memory) SelectPrev() {
	m.localStorage.SelectPrev()
}

func (m *Memory) IsEmpty() bool {
	return m.localStorage.IsEmpty()
}

func (m *Memory) GetSelectedRequest() *model.Request {
	return m.localStorage.GetSelectedRequest()
}

func (m *Memory) UpdateRequest(r *model.Request) {
	m.persistance.UpdateRequest(r)
	m.localStorage.UpdateSelectedRequest(r)
}

func (m *Memory) CreateResponse(r *model.Request) {
	m.localStorage.UpdateSelectedRequest(r)
}
