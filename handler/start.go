package handler

import (
	"database/sql"
	"log"

	"github.com/demget/quizzorobot/storage"
	tb "github.com/demget/telebot"
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
	user, err := h.db.Users.ByID(m.Sender.ID)
	if err == sql.ErrNoRows {
		err := h.db.Users.Create(m.Sender.ID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if user.State != storage.StateDefault {
		return h.sendStop(m.Sender)
	}

	_, err = h.b.Send(
		m.Sender,
		h.b.Text("start", m.Sender),
		h.b.Markup("menu"),
		tb.ModeHTML)
	return err
}

func (h Handler) onCategories(m *tb.Message) error {
	state, err := h.db.Users.State(m.Sender.ID)
	if err != nil {
		return err
	}
	if state != storage.StateDefault {
		return h.sendStop(m.Sender)
	}
	return h.sendCategories(m.Sender)
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
