package handler

import (
	"sort"
	"strconv"

	"github.com/handybots/quizzoro/bot"
	"github.com/handybots/quizzoro/storage"
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

	top = func() (filtered []storage.UserStats) {
		for _, t := range top {
			if t.Rate() > 0 {
				filtered = append(filtered, t)
			}
		}
		return
	}()

	sort.Slice(top, func(i, j int) bool {
		return top[i].Rate() > top[j].Rate()
	})

	var chats []tb.Chat
	for _, t := range top {
		chat, err := h.b.ChatByID(strconv.FormatInt(t.ID, 10))
		if err != nil {
			return err
		}
		chats = append(chats, *chat)
	}

	stats, err := h.db.Users.Stats(m.Sender.ID)
	if err != nil {
		return err
	}

	statsx := bot.Stats{
		Chats: chats,
		Top:   top,
		User:  stats,
	}

	_, err = h.b.Send(
		m.Chat,
		h.b.Text("stats", statsx),
		tb.ModeHTML)
	return err
}
