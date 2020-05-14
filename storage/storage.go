package storage

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB

	Users UsersStorage
	Polls PollsStorage
}

func Connect(url string) (*DB, error) {
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	return &DB{
		DB:    db,
		Users: &UsersTable{DB: db},
		Polls: &PollsTable{DB: db},
	}, nil
}

type Model struct {
	CreatedAt time.Time `db:"created_at" structs:"created_at"`
	UpdatedAt time.Time `db:"updated_at" structs:"updated_at"`
}
