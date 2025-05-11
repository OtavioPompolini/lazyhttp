package state

import (
	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/types"
)

type CollectionObserver interface {
	onUpdateCollection()
}

type CollectionSystem struct {
	collections []*Collection

	observers []
	persistance database.PersistanceAdapter
}

type Collection struct {
	name     string
	head     *types.Request
	tail     *types.Request
	selected *types.Request
}

func newCollectionSystem(db database.PersistanceAdapter) *CollectionSystem {
	return &CollectionSystem{
		persistance: db,
	}
}

func (c *CollectionSystem) NewCollection(collName string) {
	// nColl := &Collection{
	// 	name: collName,
	// }
	//
	// c.collections = append(c.collections, nColl)
	// c.persistance.CollectionRepository.Save(*nColl)

}

func (c *CollectionSystem) loadCollection(requests []*types.Request) {
	c.collections = append(c.collections, newCollection(requests))
}

func newCollection(requests []*types.Request) *Collection {
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
