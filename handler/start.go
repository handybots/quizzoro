package handler

import (
	"database/sql"
	"time"

	"github.com/handybots/quizzoro/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h Handler) OnStart(c tele.Context) error {
	if err := h.onStart(m); err != nil {
		h.OnError(m, err)
	}
}

func (h Handler) OnCategories(c tele.Context) error {
	return h.onCategories(c.Message())
}

func (h Handler) onStart(c tele.Context) error {
	var created bool

	user, err := h.db.Users.ByID(c.Chat().ID)
	if created = err == sql.ErrNoRows; created {
		if m.FromGroup() {
			err = h.db.Users.Create(c.Chat().ID)
		} else {
			err = h.db.Users.Create(int64(c.Sender().ID))
		}
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if user.State != storage.StateDefault {
		return h.sendStop(c.Chat())
	}

	if err := c.Send(
		h.b.Text("start", c.Chat),
		h.b.Markup("menu"),
		tb.ModeHTML,
	); err != nil {
		return err
	}

	if created && !c.Message().FromGroup() {
		<-time.After(3 * time.Second)
		return h.onSettings(c.Message())
	}
	return nil
}

func (h Handler) onCategories(c tele.Context) error {
	m := c.Message()

	state, err := h.db.Users.State(m.Chat.ID)
	if err != nil {
		return err
	}
	if state != storage.StateDefault {
		return h.sendStop(m.Chat)
	}
	return h.sendCategories(m.Chat)
}

func (h Handler) sendCategories(c tele.Context) error {
	return c.Send(
		h.b.Text("categories"),
		h.b.InlineMarkup("categories"),
		tb.ModeHTML,
	)
	return err
}

func (h Handler) sendStop(to tb.Recipient) error {
	_, err := h.b.Send(to, h.b.Text("stop"), tb.ModeHTML)
	return err
}
