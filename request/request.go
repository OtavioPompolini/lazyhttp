package request

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	REQUEST_GET    string = "get"
	REQUEST_POST          = "post"
	REQUEST_DELETE        = "delete"
	REQUEST_PUT           = "put"
)

// V1 = Only name and body
type Request struct {
	Id int64
	Name string
	// Url     string
	// Method  string
	// Headers []Header
	Body string
}

type Header struct {
	key   string
	value string
}

type SqliteDB struct {
	db *sql.DB
}

//TODO: HEHEHE Sql in model Request XD
func InitDatabase() (Adapter, error) {
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

	return sqldb, nil
}

func (a SqliteDB) GetRequests() *map[int64]Request {
	requests := make(map[int64]Request)

	row, err := a.db.Query(`
		SELECT * FROM requests
		`)
	if err != nil {
		log.Fatal(err)
	}

	defer row.Close()
	for row.Next() {
		request := Request{}

		err := row.Scan(&request.Id, &request.Name, &request.Body)
		if err != nil {
			log.Fatal(err)
		}
		requests[request.Id] = request
	}

	return &requests
}

func (a SqliteDB) CreateRequest(name string) *Request {
	res, err := a.db.Exec("INSERT INTO requests(name) values (?)", name)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return &Request{
		Id: id,
		Name: name,
		Body: "",
	}
}

func (a SqliteDB) UpdateRequest(r *Request) {
	_, err := a.db.Exec("UPDATE requests SET body=? WHERE id=?", r.Body, r.Id)
	if err != nil {
		log.Fatal(err)
	}

	// id, err := res.LastInsertId()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// return &Request{
	// 	Id: id,
	// 	Name: name,
	// 	Body: "",
	// }
}
