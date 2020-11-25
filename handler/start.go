package handler

import (
	"database/sql"
	"time"

	"github.com/handybots/quizzoro/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h Handler) OnStart(c tele.Context) error {
	return h.onStart(c)
}

func (h Handler) OnCategories(c tele.Context) error {
	return h.onCategories(c)
}

func (h Handler) onStart(c tele.Context) error {
	var created bool
	m := c.Message()

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
		return h.sendStop(c)
	}

	if err := c.Send(
		h.lt.Text(c, "start", c.Chat),
		h.lt.Markup(c, "menu"),
		tele.ModeHTML,
	); err != nil {
		return err
	}

	if created && !c.Message().FromGroup() {
		<-time.After(3 * time.Second)
		return h.onSettings(c)
	}
	return nil
}

func (h Handler) onCategories(c tele.Context) error {
	state, err := h.db.Users.State(c.Chat().ID)
	if err != nil {
		return err
	}
	if state != storage.StateDefault {
		return h.sendStop(c)
	}
	return h.sendCategories(c)
}

func (h Handler) sendCategories(c tele.Context) error {
	return c.Send(
		h.lt.Text(c, "categories"),
		h.lt.Markup(c, "categories"),
		tele.ModeHTML,
	)
}

func (h Handler) sendStop(c tele.Context) error {
	return c.Send(h.lt.Text(c, "stop"), tele.ModeHTML)
}
