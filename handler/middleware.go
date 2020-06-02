package handler

import (
	"log"
	"strings"

	"github.com/demget/quizzorobot/handler/tracker"
	tb "github.com/demget/telebot"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (h Handler) Middleware(u *tb.Update) bool {
	var f logrus.Fields

	switch {
	case u.Message != nil:
		m := u.Message
		f = eventFields(m)

		if m.Text != "" && tracker.IsSpam(m.Chat.ID, m.Text) {
			_ = h.b.Delete(m)
			f["spam"] = true
		}
	case u.Callback != nil:
		c := u.Callback
		f = eventFields(c)

		if tracker.IsSpam(c.Message.Chat.ID, c.Data) {
			f["spam"] = true
		}
	case u.PollAnswer != nil:
		f = eventFields(u.PollAnswer)
	default:
		return false
	}

	data := f["data"]
	delete(f, "data")
	logrus.WithFields(f).Info(data)

	_, spam := f["spam"]
	return !spam
}

func (h Handler) OnError(v interface{}, err error) {
	var f logrus.Fields
	if s, ok := v.(string); ok {
		f = logrus.Fields{"from": s}
	} else {
		f = eventFields(v)
	}

	logrus.WithFields(f).Error(err)
	log.Printf("%+v\n", errors.WithStack(err))
}

func eventFields(v interface{}) (f logrus.Fields) {
	var (
		user *tb.User
		chat *tb.Chat
		kind string
		data string
	)

	switch vv := v.(type) {
	case *tb.Message:
		kind = "message"
		data = vv.Text
		user = vv.Sender
		chat = vv.Chat
	case *tb.Callback:
		kind = "callback"
		data = trimData(vv.Data)
		user = vv.Sender
		chat = vv.Message.Chat
	case *tb.PollAnswer:
		kind = "poll_answer"
		data = vv.PollID
		user = &vv.User
	default:
		return
	}

	f = logrus.Fields{"event": kind}
	if data != "" {
		f["data"] = data
	}

	if user != nil {
		f["user"] = logrus.Fields{
			"id":   user.ID,
			"lang": user.LanguageCode,
		}
	}
	if chat != nil {
		f["chat"] = logrus.Fields{
			"id": chat.ID,
		}
	}

	return f
}

func trimData(s string) string {
	return strings.TrimPrefix(s, "\f")
}
