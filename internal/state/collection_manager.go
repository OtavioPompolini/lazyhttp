package state

import (
	"log"

	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/types"
)

type CollectionsChangedEvent struct {
	Collections []*types.Collection
	Cursor      int
	ActivePos   int
}

type CollectionSelectedEvent struct {
	Collection *types.Collection
}

type CollectionManager struct {
	collections []*types.Collection
	cursor      int
	activePos   int
	selectedID  int64

	eventBus             *EventBus
	collectionRepository database.CollectionRepository
}

func newCollectionManager(db database.PersistenceAdapter, eb *EventBus) *CollectionManager {
	cs := &CollectionManager{
		collectionRepository: db.CollectionRepository,
		eventBus:             eb,
	}
	return cs
}

func (c *CollectionManager) init() {
	c.collections = c.collectionRepository.GetAll()

	if len(c.collections) > 0 {
		c.selectedID = c.collections[0].Id
		c.activePos = 0
	}

	c.eventBus.Publish(c.getCollectionEvent())
	if len(c.collections) > 0 {
		c.eventBus.Publish(Event{
			Type: CollectionSelected,
			Data: CollectionSelectedEvent{Collection: c.collections[0]},
		})
	}
}

func (c *CollectionManager) getCollectionEvent() Event {
	return Event{
		Type: CollectionChanged,
		Data: CollectionsChangedEvent{
			Collections: c.collections,
			Cursor:      c.cursor,
			ActivePos:   c.activePos,
		},
	}
}

func (c *CollectionManager) NewCollection(collName string) {
	newCollection := types.Collection{
		Name:     collName,
		Position: len(c.collections),
	}

	saved := c.collectionRepository.Save(newCollection)
	c.collections = append(c.collections, saved)
	c.eventBus.Publish(c.getCollectionEvent())
}

func (c *CollectionManager) SelectNext() {
	if len(c.collections)-1 <= c.cursor {
		//Alert screen
		log.Print("Unable to select next collection. Already at the end")
		return
	}

	c.cursor += 1
	log.Println("Current collection pos:", c.cursor)
	c.eventBus.Publish(c.getCollectionEvent())
}

func (c *CollectionManager) SelectPrev() {
	if c.cursor <= 0 {
		//Alert screen
		log.Print("Unable to select previous collection. Already at the beginning")
		return
	}

	c.cursor -= 1
	log.Println("Current collection pos:", c.cursor)
	c.eventBus.Publish(c.getCollectionEvent())
}

// SWAP POSITIONS NOT WORKING CORRECTLY, I DONT CARE RN
func (c *CollectionManager) SwapPositionUp() {
	if c.cursor <= 0 {
		//Alert screen
		// Not here, this should return error and then whos calling this should
		// call alert message
		log.Print("Unable to swap collection position")
		return
	}

	c.collections[c.cursor].Position, c.collections[c.cursor-1].Position = c.collections[c.cursor-1].Position, c.collections[c.cursor].Position
	c.collections[c.cursor], c.collections[c.cursor-1] = c.collections[c.cursor-1], c.collections[c.cursor]
	c.collectionRepository.UpdatePosition(c.collections[c.cursor])
	c.collectionRepository.UpdatePosition(c.collections[c.cursor-1])

	if c.cursor == c.activePos {
		c.activePos -= 1
	} else if c.cursor-1 == c.activePos {
		c.activePos += 1
	}
	c.cursor -= 1
	c.eventBus.Publish(c.getCollectionEvent())
}

func (c *CollectionManager) SwapPositionDown() {
	if len(c.collections)-1 <= c.cursor {
		//Alert screen
		log.Print("Unable to swap collection position")
		return
	}

	c.collections[c.cursor].Position, c.collections[c.cursor+1].Position = c.collections[c.cursor+1].Position, c.collections[c.cursor].Position
	c.collections[c.cursor], c.collections[c.cursor+1] = c.collections[c.cursor+1], c.collections[c.cursor]
	c.collectionRepository.UpdatePosition(c.collections[c.cursor])
	c.collectionRepository.UpdatePosition(c.collections[c.cursor+1])

	if c.cursor == c.activePos {
		c.activePos += 1
	} else if c.cursor+1 == c.activePos {
		c.activePos -= 1
	}
	c.cursor += 1
	c.eventBus.Publish(c.getCollectionEvent())
}

func (c *CollectionManager) SelectCurrent() {
	c.activePos = c.cursor
	c.selectedID = c.collections[c.cursor].Id
	c.eventBus.Publish(c.getCollectionEvent())
}

func (c *CollectionManager) debugTestAlert() {
	e := Event{
		Type: AlertMessage,
		Data: "pudim 123",
	}

	c.eventBus.Publish(e)
}
