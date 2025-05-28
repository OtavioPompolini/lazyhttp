package state

import (
	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/types"
)

type UpdateRequestObserver interface {
	OnUpdateRequest()
}

type RequestSystem struct {
	requests             map[int][]*types.Request
	pos                  int
	selectedCollectionId int

	updateRequestObservers []UpdateRequestObserver

	requestRepository database.RequestRepository
}

func newRequestSystem(db database.PersistanceAdapter) *RequestSystem {
	return &RequestSystem{
		requestRepository: db.RequestRepository,
	}
}

func (rs *RequestSystem) SubscribeUpdateRequestEvent(obs UpdateRequestObserver) {
	rs.updateRequestObservers = append(rs.updateRequestObservers, obs)
	obs.OnUpdateRequest()
}

func (rs *RequestSystem) sendUpdateRequestEvent() {
	for _, v := range rs.updateRequestObservers {
		v.OnUpdateRequest()
	}
}

func (rs *RequestSystem) Create(reqName string) {
	saved := rs.requestRepository.Create(reqName, int64(rs.selectedCollectionId))
	currRequests, ok := rs.requests[rs.selectedCollectionId]
	if !ok {
		rs.requests[rs.selectedCollectionId] = []*types.Request{saved}
		return
	}

	currRequests = append(currRequests, saved)
	rs.pos = len(currRequests)
	rs.sendUpdateRequestEvent()
}

func (rs *RequestSystem) ListNames() []string {
	return []string{}
	// if len(rs.requests[rs.currColl]) <= 0 {
	// 	return []string{}
	// }
	//
	// requestsList := []string{}
	//
	// for _, v := range rs.requests[rs.currColl] {
	// 	requestsList = append(requestsList, v.Name)
	// }
	//
	// return requestsList
}

func (rs *RequestSystem) CurrentPos() int {
	return 0
}

// func (c *CollectionSystem) DeleteSelectedRequest() {
//
// 	selected := c.currentRequest
// 	prev := c.currentRequest.Prev
// 	next := c.currentRequest.Next
//
// 	if prev == nil && next == nil {
// 		c.currentRequest = nil
// 	}
//
// 	if prev != nil {
// 		prev.Next = selected.Next
// 		c.currentRequest = prev
// 	} else {
// 		c.currentRequest = selected.Next
// 	}
//
// 	if next != nil {
// 		next.Prev = selected.Prev
// 		c.currentRequest = next
// 	} else {
// 		c.currentRequest = selected.Prev
// 	}
//
// 	c.requestRepository.DeleteRequest(selected.Id)
// }

func (rs *RequestSystem) SelectNext() {
	currRequests, ok := rs.requests[rs.selectedCollectionId]
	if !ok {
		return
	}

	rs.pos = min(len(currRequests)-1, rs.pos+1)
}

func (rs *RequestSystem) SelectPrev() {
	_, ok := rs.requests[rs.selectedCollectionId]
	if !ok {
		return
	}

	rs.pos = max(0, rs.pos-1)
}

func (rs *RequestSystem) Update(r *types.Request) {
	// c.requests[c.currColl][c.currReq].Body = r.Body
	// c.requestRepository.UpdateRequest(r)
}

// func (c *CollectionSystem) IsEmpty() bool {
// 	return rss.state.collection.head == nil
// }

// func (rss RequestStateService) ListNames() []string {
// 	curr := rss.state.collection.head
// 	names := []string{}
//
// 	for curr != nil {
// 		names = append(names, curr.Name)
// 	}
//
// 	return names
// }

// func (rss RequestStateService) Index() int {
// 	i := 0
// 	curr := rss.state.collection.head
// 	for curr != nil {
// 		if curr.Id == rss.SelectedRequest().Id {
// 			return i
// 		}
//
// 		curr = curr.Next
// 		i += 1
// 	}
//
// 	return 0
// }

// This might not be here. XD
// func (c *CollectionSystem) ExecuteRequest() error {
// 	r := c.currentRequest
//
// 	httpRequest, err := utils.ParseHttpRequest(r.Body)
// 	if err != nil {
// 		return err
// 	}
//
// 	log.Printf("Method = %s", httpRequest.Method)
// 	log.Printf("Url = %s", httpRequest.URL)
// 	log.Printf("Body = %s", httpRequest.Body)
//
// 	client := http.Client{}
// 	res, err := client.Do(httpRequest)
// 	if err != nil {
// 		log.Print("Error while performing the request", err)
// 		return err
// 	}
//
// 	responseString := ""
//
// 	responseString += res.Proto + " "
// 	responseString += res.Status
// 	responseString += "\n"
//
// 	for k, v := range res.Header {
// 		responseString += k + ": "
// 		responseString += strings.Join(v, "")
// 		responseString += "\n"
// 	}
//
// 	responseString += "\n"
// 	s, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		log.Print("Error while reading response body", err)
// 		return err
// 	}
//
// 	// log.Printf("Response body = %s", string(s))
//
// 	response := c.responseRepository.Save(&types.Response{
// 		RequestId: c.currentRequest.Id,
// 		Info:      responseString,
// 		Body:      string(s),
// 	})
//
// 	r.ResponseHistory = append([]*types.Response{response}, r.ResponseHistory...)
//
// 	return nil
// }
