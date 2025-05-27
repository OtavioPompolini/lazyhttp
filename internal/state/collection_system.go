package state

import (
	"log"

	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/types"
)

type UpdateCollectionObserver interface {
	OnUpdateCollection()
}

type UpdateRequestObserver interface {
	OnUpdateRequest()
}

type CollectionSystem struct {
	collections []*types.Collection
	currColl    int
	selectedPos int

	requests map[int][]*types.Request
	currReq  int

	updateCollectionObservers []UpdateCollectionObserver
	updateRequestObservers    []UpdateRequestObserver

	collectionRepository database.CollectionRepository
	requestRepository    database.RequestRepository
	// responseRepository   database.ResponseRepository
}

func newCollectionSystem(db database.PersistanceAdapter) *CollectionSystem {
	return &CollectionSystem{
		// responseRepository:   db.ResponseRepository,
		requestRepository:    db.RequestRepository,
		collectionRepository: db.CollectionRepository,
		collections:          db.CollectionRepository.GetAll(),
		requests:             make(map[int][]*types.Request),
	}
}

func (c *CollectionSystem) NewCollection(collName string) {
	nColl := types.Collection{
		Name:     collName,
		Position: len(c.collections),
	}

	saved := c.collectionRepository.Save(nColl)
	c.collections = append(c.collections, saved)
	c.sendUpdateCollectionEvent()
}

func (c *CollectionSystem) SubscribeUpdateCollectionEvent(co UpdateCollectionObserver) {
	c.updateCollectionObservers = append(c.updateCollectionObservers, co)
	co.OnUpdateCollection()
}

func (c *CollectionSystem) sendUpdateCollectionEvent() {
	for _, ob := range c.updateCollectionObservers {
		ob.OnUpdateCollection()
	}
}

func (c *CollectionSystem) SelectNext() {
	if len(c.collections)-1 <= c.currColl {
		//Alert screen
		log.Print("Unable to select next collection. Already at the end")
		return
	}

	c.currColl += 1
	c.sendUpdateCollectionEvent()
}

func (c *CollectionSystem) SelectPrev() {
	if c.currColl <= 0 {
		//Alert screen
		log.Print("Unable to select previous collection. Already at the beginning")
		return
	}

	c.currColl -= 1
	c.sendUpdateCollectionEvent()
}

// SWAP POSITIONS NOT WORKING CORRECTLY, I DONT CARE RN
func (c *CollectionSystem) SwapPositionUp() {
	if c.currColl <= 1 {
		//Alert screen
		// Not here, this should return error and then whos calling this should
		// call alert message
		log.Print("Unable to swap collection position")
		return
	}

	c.collections[c.currColl].Position, c.collections[c.currColl-1].Position = c.collections[c.currColl-1].Position, c.collections[c.currColl].Position
	c.collections[c.currColl], c.collections[c.currColl-1] = c.collections[c.currColl-1], c.collections[c.currColl]
	c.collectionRepository.UpdatePosition(c.collections[c.currColl])
	c.collectionRepository.UpdatePosition(c.collections[c.currColl-1])

	if c.currColl == c.selectedPos {
		c.selectedPos -= 1
	} else if c.currColl-1 == c.selectedPos {
		c.selectedPos += 1
	}
	c.currColl -= 1

	c.sendUpdateCollectionEvent()
}

func (c *CollectionSystem) SwapPositionDown() {
	if len(c.collections)-1 <= c.currColl || c.currColl == 0 {
		//Alert screen
		log.Print("Unable to swap collection position")
		return
	}

	c.collections[c.currColl].Position, c.collections[c.currColl+1].Position = c.collections[c.currColl+1].Position, c.collections[c.currColl].Position
	c.collections[c.currColl], c.collections[c.currColl+1] = c.collections[c.currColl+1], c.collections[c.currColl]
	c.collectionRepository.UpdatePosition(c.collections[c.currColl])
	c.collectionRepository.UpdatePosition(c.collections[c.currColl+1])

	if c.currColl == c.selectedPos {
		c.selectedPos += 1
	} else if c.currColl+1 == c.selectedPos {
		c.selectedPos -= 1
	}
	c.currColl += 1

	c.sendUpdateCollectionEvent()
}

func (c *CollectionSystem) CurrentPos() int {
	return c.currColl
}

// func (c *CollectionSystem) List() []types.Collection {
// 	collectionsList := []types.Collection{}
//
// 	for _, v := range c.collections {
// 		collectionsList = append(collectionsList, *v)
// 	}
//
// 	return collectionsList
// }

func (c *CollectionSystem) SelectCurrent() {
	c.selectedPos = c.currColl
	c.sendUpdateCollectionEvent()
}

func (c *CollectionSystem) ListNames() []string {
	if len(c.collections) <= 0 {
		return []string{}
	}
	collectionsList := []string{}

	collectionsList = append(collectionsList, "*"+c.collections[c.selectedPos].Name)

	for i, v := range c.collections {
		if i != c.selectedPos {
			collectionsList = append(collectionsList, v.Name)
		}
	}

	return collectionsList
}

func (c *CollectionSystem) ListRequests() []string {
	if len(c.requests[c.currColl]) <= 0 {
		return []string{}
	}

	requestsList := []string{}

	for _, v := range c.requests[c.currColl] {
		requestsList = append(requestsList, v.Name)
	}

	return requestsList
}

func (c *CollectionSystem) CurrentRequestPosition() int {
	return c.currReq
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

func (c *CollectionSystem) SubscribeUpdateRequestEvent(obs UpdateRequestObserver) {
	c.updateRequestObservers = append(c.updateRequestObservers, obs)
	obs.OnUpdateRequest()
}

func (c *CollectionSystem) CreateRequest(reqName string) {
	saved := c.requestRepository.CreateRequest(reqName, c.collections[c.currColl].Id)
	c.requests[c.selectedPos] = append(c.requests[c.selectedPos], saved)
	c.currReq = len(c.requests[c.selectedPos])
	c.sendUpdateRequestEvent()
}

func (c *CollectionSystem) sendUpdateRequestEvent() {
	for _, v := range c.updateRequestObservers {
		v.OnUpdateRequest()
	}
}

func (c *CollectionSystem) SelectNextRequest() {
	currRequests, ok := c.requests[c.currColl]
	if !ok {
		return
	}

	c.currReq = min(len(currRequests)-1, c.currReq+1)
}

func (c *CollectionSystem) SelectPrevRequest() {
	_, ok := c.requests[c.currColl]
	if !ok {
		return
	}

	c.currReq = max(0, c.currReq-1)
}

func (c *CollectionSystem) UpdateRequest(r *types.Request) {
	c.requests[c.currColl][c.currReq].Body = r.Body
	c.requestRepository.UpdateRequest(r)
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
