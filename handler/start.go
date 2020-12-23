package handler

import (
	"github.com/handybots/quizzoro/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnStart(c tele.Context) error {
	var (
		m      = c.Message()
		group  = m.FromGroup()
		exists = false
	)

	user, err := h.db.Users.ByID(c.Chat().ID)
	if exists = err == nil; !exists {
		if group {
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

	if !group || !exists {
		if _, err := h.b.Send(
			c.Chat(),
			h.lt.Text(c, "start", c.Chat()),
			h.menuMarkup(c),
			tele.ModeHTML,
		); err != nil {
			return err
		}
	}

	if group {
		return h.OnCategories(c)
	} else if !exists {
		return h.OnSettings(c)
	}

	return nil
}

func (h handler) OnCategories(c tele.Context) error {
	state, err := h.db.Users.State(c.Chat().ID)
	if err != nil {
		return err
	}
	if state != storage.StateDefault {
		return h.sendStop(c)
	}
	return h.sendCategories(c)
}

func (h handler) sendCategories(c tele.Context) error {
	_, err := h.b.Send(
		c.Chat(),
		h.lt.Text(c, "categories"),
		h.lt.Markup(c, "categories"),
		tele.ModeHTML,
	)
	return err
}

func (h handler) sendStop(c tele.Context) error {
	_, err := h.b.Send(c.Chat(), h.lt.Text(c, "stop"), tele.ModeHTML)
	return err
}

func (h handler) menuMarkup(c tele.Context) *tele.ReplyMarkup {
	if c.Message().FromGroup() {
		return nil
	}
	return h.lt.Markup(c, "menu")
}
