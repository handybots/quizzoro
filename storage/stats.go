package storage

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
