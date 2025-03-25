package database

import (
	"database/sql"

	"github.com/OtavioPompolini/project-postman/internal/types"
)

type PersistanceAdapter struct {
	RequestRepository  RequestRepository
	ResponseRepository ResponseRepository
}

type RequestRepository interface {
	GetRequests() []*types.Request
	CreateRequest(name string) *types.Request
	UpdateRequest(r *types.Request)
	DeleteRequest(id int64)
}

type ResponseRepository interface {
	GetAll() map[int64][]*types.Response
	Save(r *types.Response) *types.Response
}

// Only sqlite for now
func NewPersistanceAdapter() (PersistanceAdapter, error) {
	db, err := sql.Open("sqlite3", "./lazycurl.db")
	if err != nil {
		return PersistanceAdapter{}, err
	}

	requestRepository := newRequestRepository(db)
	responseRepository := newResponseRepository(db)

	return PersistanceAdapter{
		RequestRepository:  requestRepository,
		ResponseRepository: responseRepository,
	}, nil
}
