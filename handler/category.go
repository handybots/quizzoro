package handler

import (
	"github.com/handybots/quizzoro/bot"
	"github.com/handybots/quizzoro/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnCategory(c tele.Context) error {
	defer h.b.Respond(c.Callback())
	return h.onCategory(c)
}

func (h handler) onCategory(c tele.Context) error {
	state, err := h.db.Users.State(c.Message().Chat.ID)
	if err != nil {
		return err
	}
	if state != storage.StateDefault {
		return nil
	}

	if err := c.Delete(); err != nil {
		return err
	}

	category := c.Data()
	if category == "random" {
		// TODO: Is not working now.

		msg, err := h.b.Send(c.Chat(), tele.Cube)
		if err != nil {
			return err
		}

		// dice value is 1-6
		// category = categoryOrder[msg.Dice.Value-1]

		r := bot.Random{
			Value:    msg.Dice.Value,
			Category: h.lt.Text(c, category),
		}

		if _, err := h.b.Send(
			c.Chat(),
			h.lt.Text(c, "random", r),
			h.lt.Markup(c, "quiz"),
			tele.ModeHTML,
		); err != nil {
			return err
		}
	} else {
		var markup *tele.ReplyMarkup
		if !c.Message().FromGroup() {
			markup = h.lt.Markup(c, "quiz")
		}

		if _, err := h.b.Send(
			c.Chat(),
			h.lt.Text(c, "chosen", category),
			markup, tele.ModeHTML,
		); err != nil {
			return err
		}
	}

	update := storage.User{State: storage.StateWaiting}
	if err := h.db.Users.Update(c.Chat().ID, update); err != nil {
		return err
	}

	return h.sendQuiz(c, category)
}
