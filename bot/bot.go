package bot

import (
	"time"

	tele "gopkg.in/tucnak/telebot.v3"
	"github.com/handybots/quizzoro/storage"
)

type Config struct {
	QuizzesChat tele.ChatID     `json:"quizzes_chat"`
	OpenPeriod  time.Duration `json:"open_period"`
}

type Random struct {
	Value    int
	Category string
}

type Stats struct {
	Chats []tele.Chat
	Top   []storage.UserStats
	User  storage.UserStats
}

var TrueFalseAnswers = []string{
	"Правда", "Ложь",
}
