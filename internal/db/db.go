package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Queries struct {
	// Add your prepared SQL statements here
}

func Initialize() (*Queries, error) {
	db, err := sql.Open("postgres", "your_connection_string")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Queries{}, nil
}
