package bot

import (
	"github.com/demget/quizzorobot/storage"
	tb "github.com/demget/telebot"
)

type Config struct {
	// TODO: Replace with tb.ChatID
	//  (tucnak/telebot@v2.2)
	QuizzesChat int64 `json:"quizzes_chat"`
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
