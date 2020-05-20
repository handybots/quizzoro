package handler

import (
	"strconv"

	"github.com/demget/quizzorobot/bot"
	tb "github.com/demget/telebot"
)

func (h Handler) OnStats(m *tb.Message) {
	if err := h.onStats(m); err != nil {
		h.OnError(m, err)
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

	var chats []tb.Chat
	for _, t := range top {
		chat, err := h.b.ChatByID(strconv.Itoa(t.ID))
		if err != nil {
			return err
		}
		chats = append(chats, *chat)
	}

	statsx := bot.Stats{
		Chats: chats,
		Top:   top,
		User:  stats,
	}

	_, err = h.b.Send(
		m.Sender,
		h.b.Text("stats", statsx),
		tb.ModeHTML)
	return err
}
