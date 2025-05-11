package database

import (
	"database/sql"
	"log"

	"github.com/OtavioPompolini/project-postman/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteResponseRepository struct {
	db *sql.DB
}

func newResponseRepository(db *sql.DB) ResponseRepository {
	return SqliteResponseRepository{
		db: db,
	}
}

func (ssr SqliteResponseRepository) Save(response *types.Response) *types.Response {
	// FOR NOW WE ALWAYS DELETE BEFORE INSERTING A NEW ONE -> In future I expect to be able to save more than one response
	_, err := ssr.db.Exec("DELETE from responses where request_id=?", response.RequestId)
	if err != nil {
		log.Panic(err)
	}

	res, err := ssr.db.Exec("INSERT INTO responses(request_id, info, body) values (?, ?, ?)", response.RequestId, response.Info, response.Body)
	if err != nil {
		log.Panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Panic(err)
	}

	return &types.Response{
		Id:        id,
		RequestId: response.RequestId,
		Body:      response.Body,
		Info:      response.Info,
	}
}

func (ssr SqliteResponseRepository) GetAll() map[int64][]*types.Response {
	responsesMap := make(map[int64][]*types.Response)

	row, err := ssr.db.Query(`
		SELECT * FROM responses
		ORDER BY created_at DESC
		`)
	if err != nil {
		log.Panic(err)
	}

	defer row.Close()
	for row.Next() {
		response := &types.Response{}

		err := row.Scan(&response.Id, &response.RequestId, &response.Info, &response.Body, &response.Created_at)
		if err != nil {
			log.Panic(err)
		}

		responsesMap[response.RequestId] = append(responsesMap[response.RequestId], response)
	}

	return responsesMap
}
