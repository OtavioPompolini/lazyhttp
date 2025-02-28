package memory

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/OtavioPompolini/project-postman/internal/model"
)

type SqliteDB struct {
	db *sql.DB
}

func initDatabase() (*SqliteDB, error) {
	db, err := sql.Open("sqlite3", "./lazycurl.db")
	if err != nil {
		return nil, err
	}

	sqldb := SqliteDB{
		db: db,
	}

	db.Exec(`
		CREATE TABLE IF NOT EXISTS requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			body TEXT NOT NULL DEFAULT ""
		);
	`)

	return &sqldb, nil
}

func (a SqliteDB) GetRequests() *map[int64]*model.Request {
	requests := make(map[int64]*model.Request)

	row, err := a.db.Query(`
		SELECT * FROM requests
		`)
	if err != nil {
		log.Panic(err)
	}

	defer row.Close()
	for row.Next() {
		request := &model.Request{}

		err := row.Scan(&request.Id, &request.Name, &request.Body)
		if err != nil {
			log.Fatal(err)
		}
		requests[request.Id] = request
	}

	return &requests
}

func (a SqliteDB) CreateRequest(name string) *model.Request {
	res, err := a.db.Exec("INSERT INTO requests(name) values (?)", name)
	if err != nil {
		log.Panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Panic(err)
	}

	return &model.Request{
		Id:   id,
		Name: name,
		Body: "",
	}
}

func (a SqliteDB) UpdateRequest(r *model.Request) {
	_, err := a.db.Exec("UPDATE requests SET body=? WHERE id=?", r.Body, r.Id)
	if err != nil {
		log.Panic(err)
	}
}
