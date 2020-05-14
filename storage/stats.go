package storage

import "go.mongodb.org/mongo-driver/bson"

type Stats struct {
	Users []User
	Stats []UserStats
}

type UserStats struct {
	Place     int
	Correct   int
	Incorrect int
}

func (s UserStats) Rate() int {
	return s.Correct * 100 / s.Incorrect
}

// statsAddFields contains pre-defined $addFields
// aggregation rules for stats-related functions.
//
// Fuck this aggregate mongo syntax...
var statsAddFields = bson.M{
	"correct": bson.M{
		"$size": bson.M{
			"$filter": bson.M{
				"input": "polls", "as": "p",
				"cond": bson.M{
					"$eq": bson.A{"$$p.correct", true},
				},
			},
		},
	},
	"incorrect": bson.M{
		"$size": bson.M{
			"$filter": bson.M{
				"input": "polls", "as": "p",
				"cond": bson.M{
					"$eq": bson.A{"$$p.correct", false},
				},
			},
		},
	},
}
