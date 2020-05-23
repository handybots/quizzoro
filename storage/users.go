package storage

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
)

type State = string

const (
	StateDefault State = "default"
	StateWaiting State = "waiting"
	StateQuiz    State = "quiz"
)

type UsersStorage interface {
	Create(id int) error
	Update(id int, user User) error
	ByID(id int) (User, error)
	State(id int) (string, error)
	Cache(id int) (UserCache, error)
	AddPoll(id int, poll PassedPoll) error
	TopStats() ([]UserStats, error)
	Stats(id int) (UserStats, error)
}

type UsersTable struct {
	*sqlx.DB
}

type User struct {
	Model     `sq:"-"`
	UserCache `sq:",flatten,omitempty"`

	ID    int    `sq:"id,omitempty"`
	State string `sq:"state,omitempty"`
}

type UserCache struct {
	LastPollID    string `sq:"last_poll_id,omitempty"`
	LastMessageID string `sq:"last_message_id,omitempty"`
	LastCategory  string `sq:"last_category,omitempty"`
}

func (db *UsersTable) Create(id int) error {
	const q = `INSERT INTO users (id) VALUES (?)`
	_, err := db.Exec(q, id)
	return err
}

func (db *UsersTable) ByID(id int) (user User, err error) {
	const q = `SELECT * FROM users WHERE id=?`
	return user, db.Get(&user, q, id)
}

func (db *UsersTable) Update(id int, user User) error {
	q, args, err := sq.
		Update("users").
		SetMap(structs.Map(user)).
		Where("id=?", id).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(q, args...)
	return err
}

func (db *UsersTable) State(id int) (state string, err error) {
	const q = `SELECT state FROM users WHERE id=?`
	return state, db.Get(&state, q, id)
}

func (db *UsersTable) Cache(id int) (cache UserCache, err error) {
	const q = `
		SELECT
		   last_poll_id,
		   last_message_id,
		   last_category
		FROM users
		WHERE id=?`

	return cache, db.Get(&cache, q, id)
}

func (db *UsersTable) AddPoll(id int, poll PassedPoll) error {
	const q = `INSERT INTO passed_polls (user_id, poll_id, correct) VALUES (?, ?, ?)`
	_, err := db.Exec(q, id, poll.PollID, poll.Correct)
	return err
}

func (db *UsersTable) TopStats() (stats []UserStats, err error) {
	const q = `
		SELECT *,
		    (
				SELECT COUNT(*) FROM passed_polls
				WHERE user_id=users.id AND correct=1
			) correct,
		    (
				SELECT COUNT(*) FROM passed_polls
				WHERE user_id=users.id AND correct=0
			) incorrect
		FROM users
		ORDER BY correct
		LIMIT 3`

	return stats, db.Select(&stats, q)
}

func (db *UsersTable) Stats(id int) (stats UserStats, err error) {
	const q = `
		SELECT *,
		    (
				SELECT COUNT(*) FROM passed_polls
				WHERE user_id=users.id AND correct=1
			) correct,
		    (
				SELECT COUNT(*) FROM passed_polls
				WHERE user_id=users.id AND correct=0
			) incorrect
		FROM users
		WHERE id=?`

	return stats, db.Get(&stats, q, id)
}
