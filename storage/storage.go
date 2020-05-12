package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	cli *mongo.Client
	*mongo.Database

	Users   *UsersColl
	Quizzes *QuizzesColl
}

func Connect(name, url string) (*DB, error) {
	opt := options.Client().ApplyURI(url)

	cli, err := mongo.Connect(context.Background(), opt)
	if err != nil {
		return nil, err
	}

	db := &DB{
		cli:      cli,
		Database: cli.Database(name),
	}

	db.Users = &UsersColl{
		Collection: db.Collection("users"),
	}
	db.Quizzes = &QuizzesColl{
		Collection: db.Collection("quizzes"),
	}

	return db, db.CreateIndexes()
}

func (db *DB) CreateIndexes() (err error) {
	_, err = db.Quizzes.Indexes().
		CreateOne(context.Background(), mongo.IndexModel{
			Keys: bson.D{
				{"category", 1},
				{"question", 1},
			},
		})
	return
}
