package app

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/types"
	"github.com/OtavioPompolini/project-postman/internal/utils"
)

type StateService struct {
	persistance database.PersistanceAdapter
	state       *State
}

type State struct {
	collection   *Collection
	variables    map[string]types.Variable
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

func loadState(db database.PersistanceAdapter) *State {
	reqs := db.RequestRepository.GetRequests()
	loadResponses(db, reqs)

	return &State{
		collection: NewCollection(reqs),
		variables:  map[string]types.Variable{},
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

func (ss StateService) DeleteSelectedRequest() {
	selected := ss.state.collection.selected
	prev := ss.state.collection.selected.Prev
	next := ss.state.collection.selected.Next

	if prev == nil && next == nil {
		ss.state.collection.selected = nil
	}

	if prev != nil {
		prev.Next = selected.Next
		ss.state.collection.selected = prev
	} else {
		ss.state.collection.head = selected.Next
	}

	if next != nil {
		next.Prev = selected.Prev
		ss.state.collection.selected = next
	} else {
		ss.state.collection.tail = selected.Prev
	}

	ss.persistance.RequestRepository.DeleteRequest(selected.Id)
}

func (ss StateService) CreateRequest(reqName string) *types.Request {
	saved := ss.persistance.RequestRepository.CreateRequest(reqName)

	if ss.state.collection.tail != nil {
		ss.state.collection.tail.Next = saved
	} else {
		ss.state.collection.head = saved
	}
	saved.Prev = ss.state.collection.tail
	ss.state.collection.selected = saved
	ss.state.collection.tail = saved

	return saved
}

func (ss *StateService) UpdateRequest(r *types.Request) {
	ss.persistance.RequestRepository.UpdateRequest(r)
	ss.state.collection.selected.Body = r.Body
}

func (ss *StateService) SelectNext() bool {
	if ss.state.collection.selected == nil {
		return false
	}

	next := ss.state.collection.selected.Next
	if next == nil {
		return false
	}

	ss.state.collection.selected = next
	return true
}

func (ss *StateService) SelectPrev() bool {
	if ss.state.collection.selected == nil {
		return false
	}
	prev := ss.state.collection.selected.Prev

	if prev == nil {
		return false
	}

	ss.state.collection.selected = prev
	return true
}

func (ss *StateService) ExecuteRequest() error {
	r := ss.state.collection.selected

	httpRequest, err := utils.ParseHttpRequest(r.Body)
	if err != nil {
		return err
	}

	log.Printf("Method = %s", httpRequest.Method)
	log.Printf("Url = %s", httpRequest.URL)
	log.Printf("Body = %s", httpRequest.Body)

	client := http.Client{}
	res, err := client.Do(httpRequest)
	if err != nil {
		log.Print("Error while performing the request", err)
		return err
	}

	responseString := ""

	responseString += res.Proto + " "
	responseString += res.Status
	responseString += "\n"

	for k, v := range res.Header {
		responseString += k + ": "
		responseString += strings.Join(v, "")
		responseString += "\n"
	}

	responseString += "\n"
	s, err := io.ReadAll(res.Body)
	if err != nil {
		log.Print("Error while reading response body", err)
		return err
	}

	log.Printf("Response body = %s", string(s))

	response := ss.persistance.ResponseRepository.Save(&types.Response{
		RequestId: ss.state.collection.selected.Id,
		Info:      responseString,
		Body:      string(s),
	})

	r.ResponseHistory = append([]*types.Response{response}, r.ResponseHistory...)

	return nil
}
