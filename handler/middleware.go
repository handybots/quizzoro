package handler

import (
	"log"
	"strings"

	tb "github.com/demget/telebot"
)

func (h Handler) Middleware(u *tb.Update) bool {
	var (
		user *tb.User
		kind string
		data string
	)

	switch {
	case u.Message != nil:
		kind = "MM"
		data = u.Message.Text
		user = u.Message.Sender
	case u.Callback != nil:
		kind = "CC"
		data = trimData(u.Callback.Data)
		user = u.Callback.Sender
	case u.PollAnswer != nil:
		kind = "PA"
		data = u.PollAnswer.PollID
		user = &u.PollAnswer.User
	default:
		return false
	}

	log.Println(kind, user.ID, data)
	return true
}

func trimData(s string) string {
	return strings.TrimPrefix(s, "\f")
}
