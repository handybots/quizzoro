package bot

import (
	"github.com/handybots/quizzoro/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

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
