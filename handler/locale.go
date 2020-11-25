package handler

import tele "gopkg.in/tucnak/telebot.v3"

func (h Handler) LocaleFunc(r tele.Recipient) string {
	// todo: implement get user lang
	// locale, _ := h.db.Users.Lang(r)
	return "ru"
}
