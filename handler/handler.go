package handler

import (
	"github.com/demget/quizzorobot/bot"
	"github.com/demget/quizzorobot/opentdb"
	"github.com/demget/quizzorobot/storage"

	tb "github.com/demget/telebot"
)

type Config struct {
	Conf bot.Config
	Bot  *tb.Bot
	DB   *storage.DB
	TDB  *opentdb.Session
}

type Handler struct {
	conf bot.Config
	b    *tb.Bot
	db   *storage.DB
	tdb  *opentdb.Session
}

func New(conf Config) Handler {
	return Handler{
		conf: conf.Conf,
		b:    conf.Bot,
		db:   conf.DB,
		tdb:  conf.TDB,
	}
}
