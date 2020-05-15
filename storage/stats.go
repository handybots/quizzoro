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

func (s Stats) Rate() int {
	return s.Correct * 100 / s.Total()
}

func (s Stats) Total() int {
	return s.Incorrect + s.Correct
}
