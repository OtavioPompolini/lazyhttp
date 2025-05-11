package database

import (
	"database/sql"
	"log"
	// _ "github.com/mattn/go-sqlite3"
)

type SqliteConfigRepository struct {
	db *sql.DB
}

func newConfigRepository(db *sql.DB) SqliteConfigRepository {
	return SqliteConfigRepository{
		db: db,
	}
}

func (cr SqliteConfigRepository) GetConfig() map[string]string {
	row, err := cr.db.Query(`
		SELECT * FROM configs
	`)
	if err != nil {
		log.Panic("Failed to initialize configs")
	}

	config := map[string]string{}

	defer row.Close()
	for row.Next() {
		k, v := "", ""

		err := row.Scan(&k, &v)
		if err != nil {
			log.Panic(err)
		}

		config[k] = v
	}

	return config
}

func (cr SqliteConfigRepository) Save(key, value string) {
	_, err := cr.db.Exec("UPDATE SET ?=?", key, value)
	if err != nil {
		log.Printf("Error while updating config key=%s value=%s", key, value)
	}
}
