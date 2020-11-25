package handler

import (
	"github.com/handybots/quizzoro/bot"
	"github.com/handybots/quizzoro/opentdb"
	"github.com/handybots/quizzoro/storage"

	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
)

type Config struct {
	Layout *layout.Layout
	Conf   bot.Config
	Bot    *tele.Bot
	DB     *storage.DB
	TDB    *opentdb.Session
}

type Handler struct {
	conf bot.Config
	b    *tele.Bot
	lt   *layout.Layout
	db   *storage.DB
	tdb  *opentdb.Session
}

func New(conf Config) Handler {
	return Handler{
		conf: conf.Conf,
		b:    conf.Bot,
		lt:   conf.Layout,
		db:   conf.DB,
		tdb:  conf.TDB,
	}
}
