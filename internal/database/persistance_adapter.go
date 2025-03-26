package database

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	"github.com/OtavioPompolini/project-postman/internal/types"
	"github.com/adrg/xdg"
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
	storagePath, err := getDBPath()
	if err != nil {
		return PersistanceAdapter{}, errors.New("Failed to create sqlite database file")
	}

	db, err := sql.Open("sqlite", storagePath)
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

func getDBPath() (string, error) {
	// Automatically handles OS-specific paths
	appDir := filepath.Join(xdg.DataHome, "YourApp")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(appDir, "lazycurl.db"), nil
}
