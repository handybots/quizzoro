package handler

import (
	"github.com/demget/quizzorobot/bot"
	"github.com/demget/quizzorobot/storage"

	tb "github.com/demget/telebot"
)

type Config struct {
	Conf bot.Config
	Bot  *tb.Bot
	DB   *storage.DB
}

type Handler struct {
	conf bot.Config
	b    *tb.Bot
	db   *storage.DB
}

func New(conf Config) Handler {
	return Handler{
		conf: conf.Conf,
		b:    conf.Bot,
		db:   conf.DB,
	}
}
