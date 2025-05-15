package database

import (
	"database/sql"

	"github.com/OtavioPompolini/project-postman/internal/types"
)

type SqliteCollectionRepository struct {
	db *sql.DB
}

func newCollectionRepository(db *sql.DB) SqliteCollectionRepository {
	return SqliteCollectionRepository{
		db: db,
	}
}

func (cr SqliteCollectionRepository) GetAll() []*types.Collection {
	return []*types.Collection{}
}

func (cr SqliteCollectionRepository) Save(c types.Collection) *types.Collection {
	return &types.Collection{
		Name: c.Name,
	}
}

func (cr SqliteCollectionRepository) SwapPositionUp(c *types.Collection) {
	return
}

func (cr SqliteCollectionRepository) SwapPositionDown(c *types.Collection) {
	return
}

// func (cr SqliteConfigRepository) GetConfig() map[string]string {
// 	row, err := cr.db.Query(`
// 		SELECT * FROM configs
// 	`)
// 	if err != nil {
// 		log.Panic("Failed to initialize configs")
// 	}
//
// 	config := map[string]string{}
//
// 	defer row.Close()
// 	for row.Next() {
// 		k, v := "", ""
//
// 		err := row.Scan(&k, &v)
// 		if err != nil {
// 			log.Panic(err)
// 		}
//
// 		config[k] = v
// 	}
//
// 	return config
// }
//
// func (cr SqliteConfigRepository) Save(key, value string) {
// 	_, err := cr.db.Exec("UPDATE SET ?=?", key, value)
// 	if err != nil {
// 		log.Printf("Error while updating config key=%s value=%s", key, value)
// 	}
// }
