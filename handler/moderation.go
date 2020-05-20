package handler

import (
	tb "github.com/demget/telebot"
)

func (h Handler) OnBadQuiz(c *tb.Callback) {
	if err := h.onBadQuiz(c); err != nil {
		h.OnError(c, err)
	}
}

func (h Handler) onBadQuiz(c *tb.Callback) error {
	if err := h.db.Polls.Delete(c.Data); err != nil {
		return err
	}
	return h.b.Delete(c.Message)
}

func (h Handler) OnBadAnswers(c *tb.Callback) {
	if err := h.onBadAnswers(c); err != nil {
		h.OnError(c, err)
	}
}

func (h Handler) onBadAnswers(c *tb.Callback) error {
	return nil
}
