package handler

import (
	"database/sql"
	"time"

	"github.com/demget/quizzorobot/storage"
	tb "github.com/demget/telebot"
)

func (h Handler) OnStart(m *tb.Message) {
	if err := h.onStart(m); err != nil {
		h.OnError(m, err)
	}
}
func (h Handler) OnSettings(m *tb.Message) {
	if err := h.onSettings(m); err != nil {
		h.OnError(m, err)
	}
}

func (h Handler) OnPrivacy(c *tb.Callback) {
	if err := h.onPrivacy(c); err != nil {
		h.OnError(c, err)
	}
}

func (h Handler) OnCategories(m *tb.Message) {
	if err := h.onCategories(m); err != nil {
		h.OnError(m, err)
	}
}

func (h Handler) onSettings(m *tb.Message) error {
	privacy, err := h.db.Users.Privacy(m.Sender.ID)
	if err != nil {
		return err
	}

	_, err = h.b.Send(
		m.Sender,
		h.b.Text("privacy"),
		h.b.InlineMarkup("privacy", privacy),
		tb.ModeHTML)
	return err
}

func (h Handler) onStart(m *tb.Message) error {
	var created bool

	user, err := h.db.Users.ByID(m.Sender.ID)
	if created = err == sql.ErrNoRows; created {
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
	if err != nil {
		return err
	}

	if created {
		<-time.After(5 * time.Second)
		return h.onSettings(m)
	}
	return nil
}

func (h Handler) onPrivacy(c *tb.Callback) error {
	defer h.b.Respond(c, &tb.CallbackResponse{
		Text: h.b.String("privacy"),
	})

	privacy, err := h.db.Users.InvertPrivacy(c.Sender.ID)
	if err != nil {
		return err
	}

	_, err = h.b.EditReplyMarkup(c.Message,
		h.b.InlineMarkup("privacy", privacy))
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
