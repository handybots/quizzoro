package storage

import (
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
)

type PollsStorage interface {
	Create(poll Poll) error
	Delete(id string) error
	ByQuestion(category, question string) (Poll, error)
	CorrectAnswer(id string) (int, error)
}

type PollsTable struct {
	*sqlx.DB
}

type Poll struct {
	Model       `structs:",omitempty"`
	ID          string  `db:"id" structs:"id,omitempty"`
	MessageID   string  `db:"message_id" structs:"message_id,omitempty"`
	ChatID      int64   `db:"chat_id" structs:"chat_id,omitempty"`
	Category    string  `db:"category" structs:"category,omitempty"`
	Difficulty  string  `db:"difficulty" structs:"difficulty,omitempty"`
	Question    string  `db:"question" structs:"question,omitempty"`
	QuestionEng string  `db:"question_eng" structs:"question_eng,omitempty"`
	Correct     string  `db:"correct" structs:"correct,omitempty"`
	CorrectEng  string  `db:"correct_eng" structs:"correct_eng,omitempty"`
	Answers     Strings `db:"answers" structs:"answers,omitempty"`
	AnswersEng  Strings `db:"answers_eng" structs:"answers_eng,omitempty"`
}

func (q Poll) MessageSig() (string, int64) {
	return q.MessageID, q.ChatID
}

type PassedPoll struct {
	Model   `structs:",omitempty"`
	ID      string
	Correct bool
}

type PassedPolls []PassedPoll

func (polls PassedPolls) Contains(pollID string) bool {
	for _, p := range polls {
		if p.ID == pollID {
			return true
		}
	}
	return false
}

func (db *PollsTable) Create(poll Poll) error {
	data := structs.Map(poll)
	data["answers"], _ = json.Marshal(data["answers"])
	data["answers_eng"], _ = json.Marshal(data["answers_eng"])

	q, args, err := sq.
		Insert("polls").
		SetMap(data).
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

func (db *PollsTable) ByQuestion(category, question string) (poll Poll, err error) {
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
