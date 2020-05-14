package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type State = string

const (
	StateDefault State = "default"
	StateQuiz    State = "quiz"
)

type UsersStorage interface {
	Create(id int) error
	Get(id int) (User, error)
	Update(id int, user User) error
	State(id int) (string, error)
	PassedPolls(id int) (UserPolls, error)
	AddPoll(id int, poll UserPoll) error
	Cache(id int) (UserCache, error)
	TopStats() (Stats, error)
	Stats(id int) (UserStats, error)
}

type UsersColl struct {
	*mongo.Collection
}

type User struct {
	ID    int        `bson:"id,omitempty"`
	State string     `bson:"state,omitempty"`
	Polls UserPolls  `bson:"polls,omitempty"`
	Cache *UserCache `bson:"cache,omitempty"`
}

type UserPoll struct {
	ID      string `bson:"id,omitempty"`
	Correct bool   `bson:"correct"`
}

type UserPolls []UserPoll

func (polls UserPolls) Contains(pollID string) bool {
	for _, p := range polls {
		if p.ID == pollID {
			return true
		}
	}
	return false
}

type UserCache struct {
	PollID    string `bson:"pollId,omitempty"`
	MessageID string `bson:"messageId,omitempty"`
	Category  string `bson:"category,omitempty"`
}

func (db *UsersColl) Create(id int) error {
	_, err := db.InsertOne(nil, User{ID: id})
	return err
}

func (db *UsersColl) Get(id int) (user User, err error) {
	return user, db.
		FindOne(nil, User{ID: id}).
		Decode(&user)
}

func (db *UsersColl) Update(id int, user User) error {
	_, err := db.UpdateOne(nil, User{ID: id}, bson.M{"$set": user})
	return err
}

func (db *UsersColl) State(id int) (s string, err error) {
	var user User
	opt := &options.FindOneOptions{
		Projection: bson.M{"state": 1},
	}

	return user.State, db.
		FindOne(nil, User{ID: id}, opt).
		Decode(&user)
}

func (db *UsersColl) PassedPolls(id int) (UserPolls, error) {
	var user User
	opt := &options.FindOneOptions{
		Projection: bson.M{"polls": 1},
	}

	return user.Polls, db.
		FindOne(nil, User{ID: id}, opt).
		Decode(&user)
}

func (db *UsersColl) AddPoll(id int, poll UserPoll) error {
	_, err := db.UpdateOne(nil, User{ID: id},
		bson.M{"$push": bson.M{"polls": poll}})
	return err
}

func (db *UsersColl) Cache(id int) (UserCache, error) {
	var user User
	opt := &options.FindOneOptions{
		Projection: bson.M{"cache": 1},
	}

	return *user.Cache, db.
		FindOne(nil, User{ID: id}, opt).
		Decode(&user)
}

func (db *UsersColl) TopStats() (Stats, error) {
	cur, err := db.Aggregate(context.Background(), bson.A{
		// bson.M{"$match": bson.M{"polls": bson.M{"$gt": "[]"}}},
		bson.M{"$addFields": statsAddFields},
		bson.M{"$sort": bson.M{"correct": -1}},
		bson.M{"$limit": 3},
	})
	if err != nil {
		return Stats{}, err
	}

	var results []struct {
		User      `bson:",inline"`
		Correct   int `bson:"correct"`
		Incorrect int `bson:"incorrect"`
	}
	if err := cur.All(context.Background(), &results); err != nil {
		return Stats{}, err
	}

	stats := make([]UserStats, len(results))
	users := make([]User, len(results))

	for i, result := range results {
		users[i] = result.User
		stats[i] = UserStats{
			Place:     i + 1,
			Correct:   result.Correct,
			Incorrect: result.Incorrect,
		}
	}

	return Stats{
		Users: users,
		Stats: stats,
	}, nil
}

func (db *UsersColl) Stats(id int) (UserStats, error) {
	cur, err := db.Aggregate(context.Background(), bson.A{
		bson.M{"$match": bson.M{"id": id}},
		bson.M{"$addFields": statsAddFields},
		bson.M{"$project": bson.M{"correct": 1, "incorrect": 1}},
	})
	if err != nil {
		return UserStats{}, err
	}

	var stats UserStats
	return stats, cur.Decode(&stats)
}
