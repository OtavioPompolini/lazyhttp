package database

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	CreateRequest(name string, collectionId int64) *types.Request
	UpdateRequest(r *types.Request)
	DeleteRequest(id int64)
}

type CollectionRepository interface {
	GetAll() []*types.Collection
	Save(c types.Collection) *types.Collection
	UpdatePosition(a *types.Collection)
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

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"sqlite3", driver)
	m.Up()

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
