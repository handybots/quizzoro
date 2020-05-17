package handler

import (
	"log"
	"strings"

	tb "github.com/demget/telebot"
	"github.com/sirupsen/logrus"
)

func (h Handler) Middleware(u *tb.Update) bool {
	var (
		user *tb.User
		kind string
		data string
	)

	switch {
	case u.Message != nil:
		kind = "message"
		data = u.Message.Text
		user = u.Message.Sender
	case u.Callback != nil:
		kind = "callback"
		data = trimData(u.Callback.Data)
		user = u.Callback.Sender
	case u.PollAnswer != nil:
		kind = "poll_answer"
		data = u.PollAnswer.PollID
		user = &u.PollAnswer.User
	default:
		return false
	}

	f := logrus.Fields{
		"event": kind,
	}
	f["user"] = logrus.Fields{
		"id":   user.ID,
		"lang": user.LanguageCode,
	}

	logrus.WithFields(f).Info(data)
	log.Println(kind, user.ID, data)

	return true
}

func trimData(s string) string {
	return strings.TrimPrefix(s, "\f")
}
