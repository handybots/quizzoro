package handler

import (
	"strconv"

	"github.com/demget/quizzorobot/bot"
	"github.com/demget/quizzorobot/storage"

	tb "github.com/demget/telebot"
)

var categories = map[string]int{
	"general":     9,
	"art":         25,
	"vehicles":    28,
	"celebrities": 26,
	"films":       11,
	"music":       12,
	"random":      -1,
}

var categoryOrder = []string{
	"general",
	"art",
	"vehicles",
	"celebrities",
	"films",
	"music",
	"random",
}

var trueFalseAnswers = []string{
	"Правда", "Ложь",
}

func (h Handler) OnCategory(c *tb.Callback) {
	defer h.b.Respond(c)
	if err := h.onCategory(c); err != nil {
		h.OnError(c, err)
	}
}

func (h Handler) onCategory(c *tb.Callback) error {
	state, err := h.db.Users.State(c.Sender.ID)
	if err != nil {
		return err
	}
	if state != storage.StateDefault {
		return nil
	}

	category := c.Data
	if category == "random" {
		msg, err := h.b.Send(c.Sender, tb.Cube)
		if err != nil {
			return err
		}
		category = categoryOrder[msg.Dice.Value-1]

		r := bot.Random{
			Value:    msg.Dice.Value,
			Category: h.b.String(category),
		}
		_, err = h.b.Edit(
			c.Message,
			h.b.Text("random", r),
			tb.ModeHTML)
		if err != nil {
			return err
		}
	} else {
		_, err := h.b.Edit(
			c.Message,
			h.b.Text("chosen", h.b.String(category)),
			tb.ModeHTML)
		if err != nil {
			return err
		}
	}

	update := storage.User{State: storage.StateWaiting}
	if err := h.db.Users.Update(c.Sender.ID, update); err != nil {
		return err
	}

	return h.sendQuiz(c.Sender, category)
}

// TODO: remove after tucnak/telebot v2.2 release
func (h Handler) forward(to tb.Recipient, m tb.Editable) (*tb.Message, error) {
	msg, chatID := m.MessageSig()
	msgID, _ := strconv.Atoi(msg)

	return h.b.Forward(to, &tb.Message{
		ID:   msgID,
		Chat: &tb.Chat{ID: chatID},
	})
}
