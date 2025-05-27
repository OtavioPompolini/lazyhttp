package database

import (
	"database/sql"
	"log"

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
	collections := []*types.Collection{}

	row, err := cr.db.Query(`
		SELECT * FROM collections
		ORDER BY position ASC
		`)
	if err != nil {
		log.Panic(err)
	}

	defer row.Close()
	for row.Next() {
		collection := &types.Collection{}

		err := row.Scan(&collection.Id, &collection.Position, &collection.Name)
		if err != nil {
			log.Panic(err)
		}

		collections = append(collections, collection)
	}

	return collections
}

func (cr SqliteCollectionRepository) Save(c types.Collection) *types.Collection {

	res, err := cr.db.Exec("INSERT INTO collections(name, position) values (?, ?)", c.Name, c.Position)
	if err != nil {
		log.Panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Panic(err)
	}

	return &types.Collection{
		Id:       id,
		Name:     c.Name,
		Position: c.Position,
	}
}

func (cr SqliteCollectionRepository) SwapPositions(a, b *types.Collection) {
	cr.db.Exec("UPDATE collections set position=? WHERE id=?", b.Position, a.Id)
	cr.db.Exec("UPDATE collections set position=? WHERE id=?", a.Position, b.Id)
}
