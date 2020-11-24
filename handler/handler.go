package handler

import (
	"github.com/handybots/quizzoro/bot"
	"github.com/handybots/quizzoro/opentdb"
	"github.com/handybots/quizzoro/storage"

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
