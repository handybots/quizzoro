package handler

import (
	"log"
	"strings"

	"github.com/demget/quizzorobot/handler/tracker"

	tb "github.com/demget/telebot"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Fields struct {
	data   string
	Event  string
	UserID int
	ChatID int64
	Spam   bool // omitempty
}

func NewFields(v interface{}) (f Fields) {
	var (
		user *tb.User
		chat *tb.Chat
	)

	switch vv := v.(type) {
	case *tb.Message:
		f.Event = "message"
		f.data = vv.Text
		user = vv.Sender
		chat = vv.Chat
	case *tb.Callback:
		f.Event = "callback"
		f.data = trimData(vv.Data)
		user = vv.Sender
		chat = vv.Message.Chat
	case *tb.PollAnswer:
		f.Event = "poll_answer"
		f.data = vv.PollID
		user = &vv.User
	default:
		return
	}

	if user != nil {
		f.UserID = user.ID
	}

	if chat != nil && chat.ID != int64(user.ID) {
		f.ChatID = chat.ID
	}

	return f
}

func (f Fields) Fields() logrus.Fields {
	return logrus.Fields{
		"event":   f.Event,
		"user_id": f.UserID,
		"chat_id": f.ChatID,
	}
}

func (h Handler) Middleware(u *tb.Update) bool {
	var f Fields

	switch {
	case u.Message != nil:
		m := u.Message
		f = NewFields(m)

		if m.Text != "" && tracker.IsSpam(m.Chat.ID, m.Text) {
			_ = h.b.Delete(m)
			f.Spam = true
		}
	case u.Callback != nil:
		c := u.Callback
		f = NewFields(c)

		if tracker.IsSpam(c.Message.Chat.ID, c.Data) {
			f.Spam = true
		}
	case u.PollAnswer != nil:
		f = NewFields(u.PollAnswer)
	default:
		return false
	}

	data := f.data
	logrus.WithFields(f.Fields()).Info(data)

	return !f.Spam
}

func (h Handler) OnError(v interface{}, err error) {
	var f Fields
	if s, ok := v.(string); ok {
		f = Fields{Event: s}
	} else {
		f = NewFields(v)
	}

	logrus.WithFields(f.Fields()).Error(err)
	log.Printf("%+v\n", errors.WithStack(err))
}

func trimData(s string) string {
	return strings.TrimPrefix(s, "\f")
}
