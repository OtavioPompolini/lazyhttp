package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/OtavioPompolini/project-postman/internal/types"
)

type SqliteRequestRepository struct {
	db *sql.DB
}

func newRequestRepository(db *sql.DB) SqliteRequestRepository {
	return SqliteRequestRepository{
		db: db,
	}
}

func (a SqliteRequestRepository) GetRequests() []*types.Request {
	requests := []*types.Request{}

	row, err := a.db.Query(`
		SELECT * FROM requests
		`)
	if err != nil {
		log.Panic(err)
	}

	defer row.Close()
	for row.Next() {
		request := &types.Request{}

		err := row.Scan(&request.Id, &request.Name, &request.Body)
		if err != nil {
			log.Panic(err)
		}

		requests = append(requests, request)
	}

	return requests
}

func (a SqliteRequestRepository) CreateRequest(name string, collectionId int64) *types.Request {
	res, err := a.db.Exec("INSERT INTO requests(name, collection_id) values (?, ?)", name, collectionId)
	if err != nil {
		log.Panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Panic(err)
	}

	return &types.Request{
		Id:   id,
		Name: name,
		Body: "",
	}
}

func (a SqliteRequestRepository) UpdateRequest(r *types.Request) {
	_, err := a.db.Exec("UPDATE requests SET body=? WHERE id=?", r.Body, r.Id)
	if err != nil {
		log.Panic(err)
	}
}

func (a SqliteRequestRepository) DeleteRequest(id int64) {
	_, err := a.db.Exec("DELETE from requests where id=?", id)
	if err != nil {
		log.Panic(err)
	}

	_, err = a.db.Exec("DELETE from responses where request_id=?", id)
	if err != nil {
		log.Panic(err)
	}
}
