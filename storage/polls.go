package storage

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
)

type PollsStorage interface {
	ByID(id string) (Poll, error)
	Create(poll Poll) error
	Delete(id string) error
	ByQuestion(category, question string) (Poll, error)
	CorrectAnswer(id string) (int, error)
	Available(userID int64, category string) (Poll, error)
}

type PollsTable struct {
	*sqlx.DB
}

type Poll struct {
	Model       `sq:"-"`
	ID          string  `sq:"id,omitempty"`
	MessageID   string  `sq:"message_id,omitempty"`
	ChatID      int64   `sq:"chat_id,omitempty"`
	Category    string  `sq:"category,omitempty"`
	Difficulty  string  `sq:"difficulty,omitempty"`
	Question    string  `sq:"question,omitempty"`
	QuestionEng string  `sq:"question_eng,omitempty"`
	Correct     string  `sq:"correct,omitempty"`
	CorrectEng  string  `sq:"correct_eng,omitempty"`
	Answers     Strings `sq:"answers,omitempty"`
	AnswersEng  Strings `sq:"answers_eng,omitempty"`
}

func (q Poll) MessageSig() (string, int64) {
	return q.MessageID, q.ChatID
}

type PassedPoll struct {
	Model   `sq:"-"`
	UserID  int    `sq:"user_id,omitempty"`
	PollID  string `sq:"poll_id,omitempty"`
	Correct bool   `sq:"correct,omitempty"`
}

type PassedPolls []PassedPoll

func (polls PassedPolls) Contains(pollID string) bool {
	for _, p := range polls {
		if p.PollID == pollID {
			return true
		}
	}
	return false
}

func (db *PollsTable) ByID(id string) (poll Poll, _ error) {
	const q = `SELECT * FROM polls WHERE id=?`
	return poll, db.Get(&poll, q, id)
}

func (db *PollsTable) Create(poll Poll) error {
	q, args, err := sq.
		Insert("polls").
		SetMap(structs.Map(poll)).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(q, args...)
	return err
}

func (db *PollsTable) Delete(id string) error {
	const q = `DELETE FROM polls WHERE id=?`
	_, err := db.Exec(q, id)
	return err
}

func (db *PollsTable) ByQuestion(category, question string) (poll Poll, _ error) {
	const q = `SELECT * FROM polls WHERE category=? AND question_eng=?`
	return poll, db.Get(&poll, q, category, question)
}

func (db *PollsTable) CorrectAnswer(id string) (int, error) {
	const q = `SELECT correct_eng, answers_eng FROM polls WHERE id=?`

	var poll Poll
	if err := db.Get(&poll, q, id); err != nil {
		return -1, err
	}

	correct := -1
	for i, a := range poll.AnswersEng {
		if a == poll.CorrectEng {
			correct = i
			break
		}
	}

	return correct, nil
}

func (db *PollsTable) Available(userID int64, category string) (poll Poll, _ error) {
	const q = `
		SELECT * FROM polls WHERE category=:category
		AND id != (SELECT orig_poll_id FROM users WHERE id=:user_id)
		AND id NOT IN (SELECT poll_id FROM passed_polls WHERE user_id=:user_id)
		ORDER BY RAND() LIMIT 1`

	stmt, err := db.PrepareNamed(q)
	if err != nil {
		return poll, err
	}

	return poll, stmt.Get(&poll, struct {
		Category string
		UserID   int64
	}{
		Category: category,
		UserID:   userID,
	})
}
