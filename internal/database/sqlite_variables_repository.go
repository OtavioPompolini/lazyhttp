package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/OtavioPompolini/project-postman/internal/types"
)

type SqliteVariablesRepository struct {
	db *sql.DB
}

func newVariablesRepository(db *sql.DB) SqliteVariablesRepository {
	db.Exec(`
		CREATE TABLE IF NOT EXISTS variables (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT NOT NULL,
			value TEXT NOT NULL DEFAULT ""
		);
	`)

	return SqliteVariablesRepository{
		db: db,
	}
}

func (a SqliteVariablesRepository) GetAll() []*types.Variable {
	variables := []*types.Variable{}

	row, err := a.db.Query(`
		SELECT * FROM variables
		`)
	if err != nil {
		log.Panic(err)
	}

	defer row.Close()
	for row.Next() {
		variable := &types.Variable{}

		err := row.Scan(&variable.Id, &variable.Key, &variable.Value)
		if err != nil {
			log.Panic(err)
		}

		variables = append(variables, variable)
	}

	return variables
}

func (a SqliteVariablesRepository) Save(key string) *types.Variable {
	res, err := a.db.Exec("INSERT INTO variables(key) values (?)", key)
	if err != nil {
		log.Panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Panic(err)
	}

	return &types.Variable{
		Id:  id,
		Key: key,
	}
}

func (a SqliteVariablesRepository) Update(variable *types.Variable) *types.Variable {
	return &types.Variable{}
}

func (a SqliteVariablesRepository) Delete(id int64) {
	_, err := a.db.Exec("DELETE from variables where id=?", id)
	if err != nil {
		log.Panic(err)
	}
}
