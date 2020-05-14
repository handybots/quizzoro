package storage

type Stats struct {
	Place     int
	Correct   int
	Incorrect int
}

type UserStats struct {
	User
	Stats
}

func (s UserStats) Rate() int {
	return s.Correct * 100 / (s.Incorrect + s.Correct)
}
