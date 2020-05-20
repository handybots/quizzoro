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
	total := s.Total()
	if total == 0 {
		return 0
	}
	return s.Correct * 100 / total
}

func (s Stats) Total() int {
	return s.Incorrect + s.Correct
}
