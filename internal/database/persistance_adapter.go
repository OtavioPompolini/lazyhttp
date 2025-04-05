package database

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	"github.com/OtavioPompolini/project-postman/internal/types"
	"github.com/adrg/xdg"
	_ "github.com/mattn/go-sqlite3"
)

type PersistanceAdapter struct {
	RequestRepository   RequestRepository
	ResponseRepository  ResponseRepository
	VariablesRepository VariablesRepository
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

type VariablesRepository interface {
	GetAll() []*types.Variable
	Save(key string) *types.Variable
	Update(*types.Variable) *types.Variable
	Delete(id int64)
}

// Only sqlite for now
func NewPersistanceAdapter() (PersistanceAdapter, error) {
	storagePath, err := getDBPath()
	if err != nil {
		return PersistanceAdapter{}, errors.New("Failed to create sqlite database file")
	}

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return PersistanceAdapter{}, err
	}

	requestRepository := newRequestRepository(db)
	responseRepository := newResponseRepository(db)
	variablesRepository := newVariablesRepository(db)

	return PersistanceAdapter{
		RequestRepository:   requestRepository,
		ResponseRepository:  responseRepository,
		VariablesRepository: variablesRepository,
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
