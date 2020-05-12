package handler

import (
	"log"

	"github.com/demget/quizzorobot/storage"
	tb "github.com/demget/telebot"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h Handler) OnStart(m *tb.Message) {
	if err := h.onStart(m); err != nil {
		log.Println(err)
	}
}

func (h Handler) OnCategories(m *tb.Message) {
	if err := h.onCategories(m); err != nil {
		log.Println(err)
	}
}

func (h Handler) onStart(m *tb.Message) error {
	user, err := h.db.Users.Get(m.Sender.ID)
	if err == mongo.ErrNoDocuments {
		err := h.db.Users.Create(m.Sender.ID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if user.State != storage.StateDefault {
		return nil // TODO: Send stop-message info
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
		return nil // TODO: Send stop-message info
	}

	_, err = h.b.Send(
		m.Sender,
		h.b.Text("categories"),
		h.b.InlineMarkup("categories"),
		tb.ModeHTML)
	return err
}
