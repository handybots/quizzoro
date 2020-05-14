package handler

import (
	"log"

	"github.com/demget/quizzorobot/bot"
	tb "github.com/demget/telebot"
)

func (h Handler) OnStats(m *tb.Message) {
	if err := h.onStats(m); err != nil {
		log.Println(err)
	}
}

func (h Handler) onStats(m *tb.Message) error {
	top, err := h.db.Users.TopStats()
	if err != nil {
		return err
	}

	stats, err := h.db.Users.Stats(m.Sender.ID)
	if err != nil {
		return err
	}

	_, err = h.b.Send(
		m.Sender,
		h.b.Text("stats", bot.Stats{Top: top, Stats: stats}),
		tb.ModeHTML)
	return err
}
