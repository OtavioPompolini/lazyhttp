package database

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	_ "github.com/mattn/go-sqlite3"

	"github.com/OtavioPompolini/project-postman/internal/types"
)

type PersistanceAdapter struct {
	RequestRepository    RequestRepository
	ResponseRepository   ResponseRepository
	VariablesRepository  VariablesRepository
	ConfigRepository     ConfigRepository
	CollectionRepository CollectionRepository
}

type RequestRepository interface {
	GetRequests() []*types.Request
	CreateRequest(name string) *types.Request
	UpdateRequest(r *types.Request)
	DeleteRequest(id int64)
}

type CollectionRepository interface {
	GetAll() []*types.Collection
	Save(c types.Collection) *types.Collection
	// Update(cName string) *types.Collection
}

type ResponseRepository interface {
	GetAll() map[int64][]*types.Response
	Save(r *types.Response) *types.Response
}

type ConfigRepository interface {
	GetConfig() map[string]string
	Update(k, v string)
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

	return PersistanceAdapter{
		RequestRepository:    newRequestRepository(db),
		ResponseRepository:   newResponseRepository(db),
		VariablesRepository:  newVariablesRepository(db),
		ConfigRepository:     newConfigRepository(db),
		CollectionRepository: newCollectionRepository(db),
	}, nil
}

func getDBPath() (string, error) {
	appDir := filepath.Join(xdg.DataHome, "LazyHttp")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(appDir, "lazyhttp.db"), nil
}
