package handler

import (
	"log"

	tb "github.com/demget/telebot"
)

func (h Handler) OnBadQuiz(c *tb.Callback) {
	if err := h.onBadQuiz(c); err != nil {
		log.Println(err)
	}
}

func (h Handler) onBadQuiz(c *tb.Callback) error {
	if err := h.db.Polls.Delete(c.Data); err != nil {
		return err
	}
	return h.b.Delete(c.Message)
}
