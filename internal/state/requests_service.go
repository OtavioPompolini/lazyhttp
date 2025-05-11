package state

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/types"
	"github.com/OtavioPompolini/project-postman/internal/utils"
)

type RequestStateService struct {
	state       *State
	persistance database.PersistanceAdapter
}

func newRequestStateService(state *State, db database.PersistanceAdapter) RequestStateService {
	return RequestStateService{
		state:       state,
		persistance: db,
	}
}

func (rss RequestStateService) DeleteSelectedRequest() {
	selected := rss.state.collection.selected
	prev := rss.state.collection.selected.Prev
	next := rss.state.collection.selected.Next

	if prev == nil && next == nil {
		rss.state.collection.selected = nil
	}

	if prev != nil {
		prev.Next = selected.Next
		rss.state.collection.selected = prev
	} else {
		rss.state.collection.head = selected.Next
	}

	if next != nil {
		next.Prev = selected.Prev
		rss.state.collection.selected = next
	} else {
		rss.state.collection.tail = selected.Prev
	}

	rss.persistance.RequestRepository.DeleteRequest(selected.Id)
}

func (rss RequestStateService) CreateRequest(reqName string) *types.Request {
	saved := rss.persistance.RequestRepository.CreateRequest(reqName)

	if rss.state.collection.tail != nil {
		rss.state.collection.tail.Next = saved
	} else {
		rss.state.collection.head = saved
	}
	saved.Prev = rss.state.collection.tail
	rss.state.collection.selected = saved
	rss.state.collection.tail = saved

	return saved
}

func (rss RequestStateService) UpdateRequest(r *types.Request) {
	rss.persistance.RequestRepository.UpdateRequest(r)
	rss.state.collection.selected.Body = r.Body
}

func (rss RequestStateService) SelectNext() bool {
	if rss.state.collection.selected == nil {
		return false
	}

	next := rss.state.collection.selected.Next
	if next == nil {
		return false
	}

	rss.state.collection.selected = next
	return true
}

func (rss RequestStateService) SelectPrev() bool {
	if rss.state.collection.selected == nil {
		return false
	}
	prev := rss.state.collection.selected.Prev

	if prev == nil {
		return false
	}

	rss.state.collection.selected = prev
	return true
}

func (rss RequestStateService) SelectedRequest() *types.Request {
	return rss.state.collection.selected
}

func (rss RequestStateService) IsEmpty() bool {
	return rss.state.collection.head == nil
}

func (rss RequestStateService) ListNames() []string {
	curr := rss.state.collection.head
	names := []string{}

	for curr != nil {
		names = append(names, curr.Name)
	}

	return names
}

func (rss RequestStateService) Index() int {
	i := 0
	curr := rss.state.collection.head
	for curr != nil {
		if curr.Id == rss.SelectedRequest().Id {
			return i
		}

		curr = curr.Next
		i += 1
	}

	return 0
}

// This might not be here. XD
func (rss RequestStateService) ExecuteRequest() error {
	r := rss.state.collection.selected

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

	// log.Printf("Response body = %s", string(s))

	response := rss.persistance.ResponseRepository.Save(&types.Response{
		RequestId: rss.state.collection.selected.Id,
		Info:      responseString,
		Body:      string(s),
	})

	r.ResponseHistory = append([]*types.Response{response}, r.ResponseHistory...)

	return nil
}
