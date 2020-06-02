package handler

import (
	"log"
	"strings"

	tb "github.com/demget/telebot"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (h Handler) Middleware(u *tb.Update) bool {
	var f logrus.Fields

	switch {
	case u.Message != nil:
		f = eventFields(u.Message)
	case u.Callback != nil:
		f = eventFields(u.Callback)
	case u.Poll != nil:
		f = eventFields(u.Poll)
	case u.PollAnswer != nil:
		f = eventFields(u.PollAnswer)
	default:
		return false
	}

	// TODO: Handle possible spam

	data := f["data"]
	delete(f, "data")
	logrus.WithFields(f).Info(data)

	return true
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
	case *tb.Poll:
		kind = "poll"
		data = vv.ID
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
