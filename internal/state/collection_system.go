package state

import (
	"log"

	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/types"
)

type CollectionEvent struct {
	Collections []*types.Collection
	CurrPos     int
	SelPos      int
}

type CollectionSelectedEvent struct {
	Collection *types.Collection
}

type CollectionSystem struct {
	collections []*types.Collection
	currPos     int
	selPos      int
	selId       int64

	eventBus             *EventBus
	collectionRepository database.CollectionRepository
}

func newCollectionSystem(db database.PersistanceAdapter, eb *EventBus) *CollectionSystem {
	cs := &CollectionSystem{
		collectionRepository: db.CollectionRepository,
		eventBus:             eb,
	}
	cs.wireEvents()
	return cs
}

func (c *CollectionSystem) wireEvents() {
	c.eventBus.Subscribe(InternalCollectionSelected, func(e Event) {
		id := e.Data.(int64)
		for i, col := range c.collections {
			if col.Id == id {
				c.selId = id
				c.selPos = i
				c.currPos = 0
				break
			}
		}
		c.eventBus.Publish(Event{
			Type: CollectionSelected,
			Data: CollectionSelectedEvent{Collection: c.collections[c.selPos]},
		})
	})
}

func (c *CollectionSystem) init() {
	c.collections = c.collectionRepository.GetAll()

	if len(c.collections) > 0 {
		c.selId = c.collections[0].Id
		c.selPos = 0
	}
	c.eventBus.Publish(c.getCollectionEvent())

	if len(c.collections) > 0 {
		c.eventBus.Publish(Event{
			Type: CollectionSelected,
			Data: CollectionSelectedEvent{Collection: c.collections[0]},
		})
	}
}

func (c *CollectionSystem) getCollectionEvent() Event {
	return Event{
		Type: CollectionChanged,
		Data: CollectionEvent{
			Collections: c.collections,
			CurrPos:     c.currPos,
			SelPos:      c.selPos,
		},
	}
}

func (c *CollectionSystem) NewCollection(collName string) {
	nColl := types.Collection{
		Name:     collName,
		Position: len(c.collections),
	}

	saved := c.collectionRepository.Save(nColl)
	c.collections = append(c.collections, saved)
	c.eventBus.Publish(c.getCollectionEvent())
}

func (c *CollectionSystem) SelectNext() {
	if len(c.collections)-1 <= c.currPos {
		//Alert screen
		log.Print("Unable to select next collection. Already at the end")
		return
	}

	c.currPos += 1
	log.Println("Current collection pos:", c.currPos)
}

func (c *CollectionSystem) SelectPrev() {
	if c.currPos <= 0 {
		//Alert screen
		log.Print("Unable to select previous collection. Already at the beginning")
		return
	}

	c.currPos -= 1
	log.Println("Current collection pos:", c.currPos)
}

// SWAP POSITIONS NOT WORKING CORRECTLY, I DONT CARE RN
func (c *CollectionSystem) SwapPositionUp() {
	if c.currPos <= 1 {
		//Alert screen
		// Not here, this should return error and then whos calling this should
		// call alert message
		log.Print("Unable to swap collection position")
		return
	}

	c.collections[c.currPos].Position, c.collections[c.currPos-1].Position = c.collections[c.currPos-1].Position, c.collections[c.currPos].Position
	c.collections[c.currPos], c.collections[c.currPos-1] = c.collections[c.currPos-1], c.collections[c.currPos]
	c.collectionRepository.UpdatePosition(c.collections[c.currPos])
	c.collectionRepository.UpdatePosition(c.collections[c.currPos-1])

	if c.currPos == c.selPos {
		c.selPos -= 1
	} else if c.currPos-1 == c.selPos {
		c.selPos += 1
	}
	c.currPos -= 1
	c.eventBus.Publish(c.getCollectionEvent())
}

func (c *CollectionSystem) SwapPositionDown() {
	if len(c.collections)-1 <= c.currPos || c.currPos == 0 {
		//Alert screen
		log.Print("Unable to swap collection position")
		return
	}

	c.collections[c.currPos].Position, c.collections[c.currPos+1].Position = c.collections[c.currPos+1].Position, c.collections[c.currPos].Position
	c.collections[c.currPos], c.collections[c.currPos+1] = c.collections[c.currPos+1], c.collections[c.currPos]
	c.collectionRepository.UpdatePosition(c.collections[c.currPos])
	c.collectionRepository.UpdatePosition(c.collections[c.currPos+1])

	if c.currPos == c.selPos {
		c.selPos += 1
	} else if c.currPos+1 == c.selPos {
		c.selPos -= 1
	}
	c.currPos += 1
	c.eventBus.Publish(c.getCollectionEvent())
}

func (c *CollectionSystem) SelectCurrent() {
	c.selPos = c.currPos
	c.selId = c.collections[c.currPos].Id
	c.currPos = 0
}
