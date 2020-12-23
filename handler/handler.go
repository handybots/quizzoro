package handler

import (
	"github.com/handybots/quizzoro/opentdb"
	"github.com/handybots/quizzoro/storage"

	tele "gopkg.in/tucnak/telebot.v3"
	"gopkg.in/tucnak/telebot.v3/layout"
)

func New(h Handler) handler {
	return handler{
		b:   h.Bot,
		lt:  h.Layout,
		db:  h.DB,
		tdb: h.TDB,
	}
}

type (
	Handler struct {
		Layout *layout.Layout
		Bot    *tele.Bot
		DB     *storage.DB
		TDB    *opentdb.Session
	}

	handler struct {
		b   *tele.Bot
		lt  *layout.Layout
		db  *storage.DB
		tdb *opentdb.Session
	}
)
