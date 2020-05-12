package storage

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QuizzesStorage interface {
	Create(q Quiz) error
	Delete(id string) error
	ByQuestion(q string) (Quiz, error)
	CorrectAnswer(id string) (int, error)
}

type QuizzesColl struct {
	*mongo.Collection
}

type Quiz struct {
	PollID     string   `bson:"pollId,omitempty"`
	MessageID  string   `bson:"messageId,omitempty"`
	ChatID     int64    `bson:"chatId,omitempty"`
	Category   string   `bson:"category,omitempty"`
	Difficulty string   `bson:"difficulty,omitempty"`
	Question   string   `bson:"question,omitempty"`
	Correct    string   `bson:"correct,omitempty"`
	Answers    []string `bson:"answers,omitempty"`
}

func (q Quiz) MessageSig() (string, int64) {
	return q.MessageID, q.ChatID
}

func (db *QuizzesColl) Create(q Quiz) error {
	_, err := db.InsertOne(nil, q)
	return err
}

func (db *QuizzesColl) Delete(id string) error {
	_, err := db.DeleteOne(nil, Quiz{PollID: id})
	return err
}

func (db *QuizzesColl) ByQuestion(cat, q string) (quiz Quiz, err error) {
	return quiz, db.
		FindOne(nil, Quiz{
			Category: cat,
			Question: q,
		}).
		Decode(&quiz)
}

func (db *QuizzesColl) CorrectAnswer(id string) (int, error) {
	var quiz Quiz
	opt := &options.FindOneOptions{
		Projection: bson.M{"correct": 1, "answers": 1},
	}

	err := db.FindOne(nil, Quiz{PollID: id}, opt).
		Decode(&quiz)
	if err != nil {
		return -1, err
	}

	correct := -1
	for i, a := range quiz.Answers {
		if a == quiz.Correct {
			correct = i
			break
		}
	}

	return correct, nil
}
