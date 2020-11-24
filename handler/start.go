package handler

import (
	"database/sql"
	"time"

	tb "github.com/demget/telebot"
	"github.com/handybots/quizzoro/storage"
)

func (h Handler) OnStart(m *tb.Message) {
	if err := h.onStart(m); err != nil {
		h.OnError(m, err)
	}
}

func (h Handler) OnCategories(m *tb.Message) {
	if err := h.onCategories(m); err != nil {
		h.OnError(m, err)
	}
}

func (h Handler) onStart(m *tb.Message) error {
	var created bool

	user, err := h.db.Users.ByID(m.Chat.ID)
	if created = err == sql.ErrNoRows; created {
		if m.FromGroup() {
			err = h.db.Users.Create(m.Chat.ID)
		} else {
			err = h.db.Users.Create(int64(m.Sender.ID))
		}
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if user.State != storage.StateDefault {
		return h.sendStop(m.Chat)
	}

	_, err = h.b.Send(
		m.Chat,
		h.b.Text("start", m.Chat),
		h.b.Markup("menu"),
		tb.ModeHTML)
	if err != nil {
		return err
	}

	if created && !m.FromGroup() {
		<-time.After(3 * time.Second)
		return h.onSettings(m)
	}
	return nil
}

func (h Handler) onCategories(m *tb.Message) error {
	state, err := h.db.Users.State(m.Chat.ID)
	if err != nil {
		return err
	}
	if state != storage.StateDefault {
		return h.sendStop(m.Chat)
	}
	return h.sendCategories(m.Chat)
}

func (h Handler) sendCategories(to tb.Recipient) error {
	_, err := h.b.Send(to,
		h.b.Text("categories"),
		h.b.InlineMarkup("categories"),
		tb.ModeHTML)
	return err
}

func (h Handler) sendStop(to tb.Recipient) error {
	_, err := h.b.Send(to, h.b.Text("stop"), tb.ModeHTML)
	return err
}
