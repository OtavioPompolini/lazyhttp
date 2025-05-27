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
	OnChangeSelectedCollection()
}

type CollectionSystem struct {
	collections []*types.Collection
	currPos     int
	selectedPos int

	requests map[int][]*types.Request
	currReq  int

	observers            []UpdateCollectionObserver
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
	}
}

func (c *CollectionSystem) NewCollection(collName string) {
	nColl := types.Collection{
		Name:     collName,
		Position: len(c.collections),
	}

	saved := c.collectionRepository.Save(nColl)
	c.collections = append(c.collections, saved)
	c.sendEvent()
}

func (c *CollectionSystem) Subscribe(co UpdateCollectionObserver) {
	c.observers = append(c.observers, co)
}

func (c *CollectionSystem) sendEvent() {
	for _, ob := range c.observers {
		ob.OnUpdateCollection()
	}
}

func (c *CollectionSystem) SelectNext() {
	if len(c.collections)-1 <= c.currPos {
		//Alert screen
		log.Print("Unable to select next collection. Already at the end")
		return
	}

	c.currPos += 1
	c.sendEvent()
}

func (c *CollectionSystem) SelectPrev() {
	if c.currPos <= 0 {
		//Alert screen
		log.Print("Unable to select previous collection. Already at the beginning")
		return
	}

	c.currPos -= 1
	c.sendEvent()
}

func (c *CollectionSystem) SwapPositionUp() {
	if c.currPos <= 1 {
		//Alert screen
		log.Print("Unable to swap collection position")
		return
	}

	c.collectionRepository.SwapPositions(c.collections[c.currPos], c.collections[c.currPos-1])
	c.collections[c.currPos], c.collections[c.currPos-1] = c.collections[c.currPos-1], c.collections[c.currPos]
	if c.currPos == c.selectedPos {
		c.selectedPos -= 1
	} else if c.currPos-1 == c.selectedPos {
		c.selectedPos += 1
	}
	c.currPos -= 1

	c.sendEvent()
}

func (c *CollectionSystem) SwapPositionDown() {
	if len(c.collections)-1 <= c.currPos || c.currPos == 0 {
		//Alert screen
		log.Print("Unable to swap collection position")
		return
	}

	c.collectionRepository.SwapPositions(c.collections[c.currPos], c.collections[c.currPos+1])
	c.collections[c.currPos], c.collections[c.currPos+1] = c.collections[c.currPos+1], c.collections[c.currPos]
	if c.currPos == c.selectedPos {
		c.selectedPos += 1
	} else if c.currPos+1 == c.selectedPos {
		c.selectedPos -= 1
	}
	c.currPos += 1

	c.sendEvent()
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
	c.selectedPos = c.currPos
	c.sendEvent()
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

func (c *CollectionSystem) CurrentPos() int {
	return c.currPos
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

// func (c *CollectionSystem) CreateRequest(reqName string) *types.Request {
// 	saved := c.requestRepository.CreateRequest(reqName)
//
// 	if rss.state.collection.tail != nil {
// 		rss.state.collection.tail.Next = saved
// 	} else {
// 		rss.state.collection.head = saved
// 	}
// 	saved.Prev = rss.state.collection.tail
// 	rss.state.collection.selected = saved
// 	rss.state.collection.tail = saved
//
// 	return saved
// }
//
// func (rss RequestStateService) UpdateRequest(r *types.Request) {
// 	rss.persistance.RequestRepository.UpdateRequest(r)
// 	rss.state.collection.selected.Body = r.Body
// }
//
// func (rss RequestStateService) SelectNext() bool {
// 	if rss.state.collection.selected == nil {
// 		return false
// 	}
//
// 	next := rss.state.collection.selected.Next
// 	if next == nil {
// 		return false
// 	}
//
// 	rss.state.collection.selected = next
// 	return true
// }

// func (rss RequestStateService) SelectPrev() bool {
// 	if rss.state.collection.selected == nil {
// 		return false
// 	}
// 	prev := rss.state.collection.selected.Prev
//
// 	if prev == nil {
// 		return false
// 	}
//
// 	rss.state.collection.selected = prev
// 	return true
// }

// func (c *CollectionSystem) SelectedRequest() *types.Request {
// 	return c.currentRequest
// }

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
