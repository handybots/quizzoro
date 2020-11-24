package handler

import (
	"github.com/handybots/quizzoro/bot"
	"github.com/handybots/quizzoro/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

var categories = map[string][]int{
	"general":   {9},
	"history":   {23},
	"music":     {12},
	"books":     {10},
	"games":     {15, 16},
	"computers": {18, 30},
	"random":    {-1},
}

var categoryOrder = []string{
	"general",
	"history",
	"music",
	"books",
	"games",
	"computers",
	"random",
}

func (h Handler) OnCategory(c tele.Context) error {
	defer h.b.Respond(c.Callback())
	return h.onCategory(c)
}

func (h Handler) onCategory(c tele.Context) error {
	state, err := h.db.Users.State(c.Message().Chat.ID)
	if err != nil {
		return err
	}
	if state != storage.StateDefault {
		return nil
	}

	_ = c.Delete()

	category := c.Data()
	if category == "random" {
		msg, err := h.b.Send(c.Message().Chat, tele.Cube)
		if err != nil {
			return err
		}

		// dice value is 1-6
		category = categoryOrder[msg.Dice.Value-1]

		r := bot.Random{
			Value:    msg.Dice.Value,
			Category: h.b.String(category),
		}
		if err := c.Send(
			h.b.Text("random", r),
			h.b.Markup("quiz"),
			tele.ModeHTML,
		); err != nil {
			return err
		}
	} else {
		if err := c.Send(
			h.b.Text("chosen", h.b.String(category)),
			h.b.Markup("quiz"),
			tb.ModeHTML,
		); err != nil {
			return err
		}
	}

	update := storage.User{State: storage.StateWaiting}
	if err := h.db.Users.Update(c.Message.Chat.ID, update); err != nil {
		return err
	}

	return h.sendQuiz(c.Message.Chat, category)
}
