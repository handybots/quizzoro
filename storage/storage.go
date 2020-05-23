package storage

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
	"unicode"

	"github.com/fatih/structs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

func init() {
	// structs is used with squirrel (sq)
	structs.DefaultTagName = "sq"
}

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

	db.Mapper = reflectx.NewMapperFunc("db", toSnakeCase)

	return &DB{
		DB:    db,
		Users: &UsersTable{DB: db},
		Polls: &PollsTable{DB: db},
	}, nil
}

func toSnakeCase(s string) string {
	runes := []rune(s)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) {
			if (i+1 < length && unicode.IsLower(runes[i+1])) ||
				unicode.IsLower(runes[i-1]) {
				out = append(out, '_')
			}
		}

		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

type Model struct {
	CreatedAt time.Time `sq:"created_at"`
	UpdatedAt time.Time `sq:"updated_at"`
}

type Strings []string

func (s Strings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Strings) Scan(src interface{}) error {
	if v, ok := src.([]uint8); ok {
		return json.Unmarshal(v, &s)
	}
	return errors.New("storage: Strings must be used with json field only")
}
