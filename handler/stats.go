package handler

import (
	"sort"
	"strconv"

	"github.com/handybots/quizzoro/bot"
	"github.com/handybots/quizzoro/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnStats(c tele.Context) error {
	return h.onStats(c)
}

func (h handler) onStats(c tele.Context) error {
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

	var chats []tele.Chat
	for _, t := range top {
		chat, err := h.b.ChatByID(strconv.FormatInt(t.ID, 10))
		if err != nil {
			return err
		}
		chats = append(chats, *chat)
	}

	stats, err := h.db.Users.Stats(c.Sender().ID)
	if err != nil {
		return err
	}

	statsx := bot.Stats{
		Chats: chats,
		Top:   top,
		User:  stats,
	}

	return c.Send(
		h.lt.Text(c, "stats", statsx),
		tele.ModeHTML,
	)
}
