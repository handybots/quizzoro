package handler

import (
	"github.com/demget/quizzorobot/bot"
	"github.com/demget/quizzorobot/storage"
	tb "github.com/demget/telebot"
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

func (h Handler) OnCategory(c *tb.Callback) {
	defer h.b.Respond(c)
	if err := h.onCategory(c); err != nil {
		h.OnError(c, err)
	}
}

func (h Handler) onCategory(c *tb.Callback) error {
	state, err := h.db.Users.State(c.Message.Chat.ID)
	if err != nil {
		return err
	}
	if state != storage.StateDefault {
		return nil
	}

	_ = h.b.Delete(c.Message)

	category := c.Data
	if category == "random" {
		msg, err := h.b.Send(c.Message.Chat, tb.Cube)
		if err != nil {
			return err
		}

		// dice value is 1-6
		category = categoryOrder[msg.Dice.Value-1]

		r := bot.Random{
			Value:    msg.Dice.Value,
			Category: h.b.String(category),
		}
		_, err = h.b.Send(
			c.Message.Chat,
			h.b.Text("random", r),
			h.b.Markup("quiz"),
			tb.ModeHTML)
		if err != nil {
			return err
		}
	} else {
		_, err := h.b.Send(
			c.Message.Chat,
			h.b.Text("chosen", h.b.String(category)),
			h.b.Markup("quiz"),
			tb.ModeHTML)
		if err != nil {
			return err
		}
	}

	update := storage.User{State: storage.StateWaiting}
	if err := h.db.Users.Update(c.Message.Chat.ID, update); err != nil {
		return err
	}

	return h.sendQuiz(c.Message.Chat, category)
}
