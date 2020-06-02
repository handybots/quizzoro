package bot

import (
	"time"

	"github.com/demget/quizzorobot/storage"
	tb "github.com/demget/telebot"
)

type Config struct {
	QuizzesChat tb.ChatID     `json:"quizzes_chat"`
	OpenPeriod  time.Duration `json:"open_period"`
}

type Random struct {
	Value    int
	Category string
}

type Stats struct {
	Chats []tb.Chat
	Top   []storage.UserStats
	User  storage.UserStats
}

var TrueFalseAnswers = []string{
	"Правда", "Ложь",
}
